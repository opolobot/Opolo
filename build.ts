const isProd = Deno.args[0] === "--prod"

const utilsPkgLoc = "github.com/zorbyte/whiskey/utils"
const tagsCmd = ["git", "describe", "--tags", "--long", "`git rev-list --tags --max-count=1"]

if (isProd) {
    console.log("building for production")
}

const decoder = new TextDecoder()

await checkExecutable("git")
await checkExecutable("go")

const ver = Deno.run({
    stdout: "piped",
    stderr: "piped",
    cmd: tagsCmd
})

let versionTag = "unknown"

const output = await Promise.race([ver.output(), ver.stderrOutput()])
const verOutput = decoder.decode(output)
if (!verOutput.includes("fatal")) versionTag = verOutput.slice(1)

const env = Deno.env.toObject()
const username = env.USER ?? env.UserName

const buildCmd = [
    "go",
    isProd ? "build" : "run",
    "-ldflags",
    `"${[
        `Version=${versionTag}`,
        `BuiltOn=${Date.now()}`,
        `BuiltBy=${username}`
    ].map(item => `-X ${utilsPkgLoc}.` + item).join(" ")}"`
]

if (isProd) buildCmd.push("-i")
buildCmd.push(".")

console.log(`running:\n${buildCmd.join(" ")}`)

const goBuild = Deno.run({
    cwd: Deno.cwd(),
    cmd: buildCmd
})

const stats = await goBuild.status()
Deno.exit(stats.code)

async function checkExecutable(name: string) {
    const execExists = Deno.run({
        stdout: "piped",
        stderr: "piped",
        cmd: ["which", name]
    })

    if (!(await execExists.status()).success) {
        const out = await execExists.stderrOutput()
        if (decoder.decode(out).includes("not found")) {
            console.log(`install ${name} before continuing`)
            Deno.exit(1)
        }
    }
}
