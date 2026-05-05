# confmerge

> A tool for deep-merging layered YAML/TOML config files with override tracking and diff output.

---

## Installation

```bash
go install github.com/yourname/confmerge@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/confmerge.git && cd confmerge && go build -o confmerge .
```

---

## Usage

Merge multiple config files in order (later files override earlier ones):

```bash
confmerge base.yaml staging.yaml overrides.yaml
```

Output a diff showing which keys were overridden:

```bash
confmerge --diff base.toml prod.toml
```

Write the merged result to a file:

```bash
confmerge --out merged.yaml base.yaml env.yaml secrets.yaml
```

### Example

**base.yaml**
```yaml
server:
  host: localhost
  port: 8080
database:
  pool: 5
```

**prod.yaml**
```yaml
server:
  host: 0.0.0.0
database:
  pool: 20
```

```bash
confmerge --diff base.yaml prod.yaml
```

```
~ server.host: "localhost" → "0.0.0.0"
~ database.pool: 5 → 20
```

---

## Supported Formats

- YAML (`.yaml`, `.yml`)
- TOML (`.toml`)

---

## License

MIT © yourname