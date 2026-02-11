# envsync

Securely share .env files across your team. No cloud service needed.

## Install

```bash
git clone https://github.com/openkickstartai/envsync.git
cd envsync && go build -o envsync .
```

## Usage

```bash
# Initialize in your project
envsync init

# Encrypt and commit .env for team
envsync push .env --env production

# Pull and decrypt
envsync pull --env production

# Add a team member
envsync add-key teammate@company.com public.key

# Show diff between local and encrypted
envsync diff --env staging
```

## How it works

1. `envsync init` creates `.envsync/` directory and generates an age keypair
2. `envsync push` encrypts `.env` with all team members' public keys, stores in `.envsync/`
3. `.envsync/` is committed to git (encrypted), `.env` stays in `.gitignore`
4. `envsync pull` decrypts using your private key

## Testing

```bash
go test -v ./...
```
