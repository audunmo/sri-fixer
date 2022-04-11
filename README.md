# SRI Fixer

SRI Fixer finds HTML files in your project, and automatically adds Subresource Integrity hashes to `<script>` and `<link>` tags

## Requirements

Installing SRI Fixer requires that your system has [Go](https://go.dev) installed

## Installation

We don't currently publish prebuilt binaries, so you'll need to compile sri-fixer for your system:

```sh
git clone https://github.com/artlovecode/sri-fixer
cd sri-fixer
go install

# If all has gone well, this should print a message from sri-fixer
sri-fixer

# If your shell doesn't find sri-fixer, ensure that $(go env GOPATH) is in your $PATH
```

## Usage

SRI Fixer is designed with the philosophy that it should just Do The Right Thing™, and do only _one thing_, without config or input
Navigate to your project root, and run `sri-fixer run --host "https://your-sites-domain.xyz"`. It will update the HTML files in-place

### Roadmap:
- File finder should respect .gitignore
- Publish to brew
- Debug logs
- Automated testing in testcafe
- Optional hashing of self-hosted `<scripts>/<links>`
- Ignoring multiple domains
