# opolo

opolo, the cutest opossum on Discord!

## Structure

```
├───assets
├───db
├───ocl - opolo command library
│   ├───args - Argument parsing
│   ├───embeds - UI/UX theming with embeds
│   └───msgcol - Message collector
├───pieces
│   ├───cmds - Commands
│   │   └───[category] - Commands are grouped into their respective categories by folder
│   ├───events
│   ├───mdware - Middleware
│   ├───models - Database models
│   └───parsers - Argument Parsers
└───utils - Utilities
```