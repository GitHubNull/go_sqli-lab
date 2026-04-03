# Agent Instructions

## Critical Rules for AI Agents

### READ-ONLY Dependencies

#### ref/sqli-labs (Git Submodule)

**⚠️ CRITICAL: This directory is a READ-ONLY dependency**

- **Path**: `ref/sqli-labs/`
- **Type**: Git submodule (external dependency)
- **Source**: https://github.com/Audi-1/sqli-labs

**FORBIDDEN ACTIONS**:
- ❌ NEVER modify any files in `ref/sqli-labs/`
- ❌ NEVER delete files in `ref/sqli-labs/`
- ❌ NEVER add new files to `ref/sqli-labs/`
- ❌ NEVER rename files in `ref/sqli-labs/`

**ALLOWED ACTIONS**:
- ✅ Read and reference files for implementation
- ✅ Study the PHP code to understand vulnerability patterns
- ✅ Use as reference when porting to Go

**HOW TO UPDATE**:
Only sync updates when the upstream repository changes:

```bash
# Fetch latest changes from upstream
git submodule update --init --remote ref/sqli-labs

# Or update all submodules
git submodule update --init --remote --recursive
```

**WHY THIS MATTERS**:
This is the reference PHP implementation of sqli-labs. We are porting it to Go in the `src/` directory. Keeping it as a git submodule ensures:
1. We can always reference the original implementation
2. We can track upstream updates easily
3. We maintain a clean separation between our code and third-party code

---

## Project Structure Rules

1. **Go source code**: All implementation in `src/` directory
2. **Static resources**: All web assets in `src/web/`
3. **Vulnerability implementations**: All in `src/vulnerabilities/` as .go files
4. **Configuration**: All in `config/` directory
5. **Documentation**: All in `doc/` directory

## Development Guidelines

- Follow existing code patterns in the codebase
- Each vulnerability should implement the `Vulnerability` interface
- Use the registry pattern to register new vulnerabilities
- Maintain compatibility with multiple database types (SQLite, MySQL, PostgreSQL)
- Never suppress type errors with `as any` or `@ts-ignore`
