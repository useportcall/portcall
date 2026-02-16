.PHONY: dev-cli test-all test-go-workspace test-go-libs-non-services test-integration test-dashboard-fe test-email e2e-email-local e2e-email-live test-full e2e e2e-verbose e2e-browser e2e-browser-headed e2e-discord e2e-browser-discord e2e-discord-live e2e-browser-discord-live e2e-browser-snapshots e2e-browser-snapshots-live

# Build the dev-cli and copy to project root
dev-cli:
	cd tools/dev-cli && go build -o ../../dev-cli .

# Run all Go tests for every module listed in go.work (includes e2e package).
test-go-workspace:
	@set -euo pipefail; \
	mods=$$(awk '/^\t\.\//{gsub(/^\t/, ""); print $$1}' go.work); \
	for m in $$mods; do \
		echo "==> $$m"; \
		(cd $$m && go test -count=1 ./...); \
	done

# Run go test and go vet for libs/go modules except libs/go/services.
test-go-libs-non-services:
	@set -euo pipefail; \
	mods=$$(find libs/go -mindepth 1 -maxdepth 1 -type d | sort); \
	for m in $$mods; do \
		name=$$(basename $$m); \
		if [ "$$name" = "services" ]; then continue; fi; \
		echo "==> $$m (test)"; \
		(cd $$m && go test -count=1 ./...); \
		echo "==> $$m (vet)"; \
		(cd $$m && go vet ./...); \
	done

# Run integration-tagged tests.
test-integration:
	go test -count=1 -tags=integration ./libs/go/services/...

# Run dashboard frontend tests.
test-dashboard-fe:
	cd apps/dashboard/frontend && npm test

# Run focused email-related tests.
test-email:
	go test -count=1 ./libs/go/envx/... ./libs/go/emailx/... ./apps/email/...

# Run local end-to-end email tests with a mock Resend API.
e2e-email-local:
	go test -count=1 -tags=e2e_email_local -run TestEmailWorkerLocalE2E ./apps/email/...
	go test -count=1 -tags=e2e_email_local -run TestSMTPRelayLocalE2E ./apps/email/internal/smtprelay/...

# Run live end-to-end email tests against real Resend (sends real emails).
e2e-email-live:
	@test -n "$$RESEND_API_KEY" && test -n "$$E2E_EMAIL_FROM" && test -n "$$E2E_EMAIL_TO" || (echo "Set RESEND_API_KEY, E2E_EMAIL_FROM, and E2E_EMAIL_TO"; exit 1)
	go test -count=1 -tags=e2e_live_email -run TestEmailWorkerLiveResend ./apps/email/...
	go test -count=1 -tags=e2e_live_email -run TestSMTPRelayLiveResend ./apps/email/internal/smtprelay/...

# Run all tests: Go workspace, integration, email, and dashboard frontend.
test-full: test-go-workspace test-integration test-email test-dashboard-fe

# Run all cross-app e2e tests (auto-starts Postgres if needed)
e2e:
	go test -count=1 -timeout 120s ./e2etest/

# Run all Go tests in the workspace (including e2e package tests)
test-all:
	go test -count=1 $(shell go list -m -f '{{if .Main}}{{.Dir}}/...{{end}}' all)

# Run e2e tests with verbose output
e2e-verbose:
	go test -v -count=1 -timeout 120s ./e2etest/

# Run browser-based e2e tests (headless)
e2e-browser:
	cd e2etest/browser && npx playwright test

# Run browser-based e2e tests (visible browser)
e2e-browser-headed:
	cd e2etest/browser && npx playwright test --headed

# Run focused Discord notification e2e test (captures webhook posts locally).
e2e-discord:
	go test -count=1 -timeout 120s ./e2etest/ -run TestE2E_DiscordNotifications

# Run focused Discord notification browser test (captures webhook posts locally).
e2e-browser-discord:
	cd e2etest/browser && npx playwright test -g "discord notifications fire"

# Run Discord notification e2e against real webhook URLs loaded from app .envs files.
e2e-discord-live:
	@set -euo pipefail; \
	for f in apps/dashboard/.envs apps/api/.envs apps/billing/.envs; do \
		if [ -f "$$f" ]; then set -a; . "$$f"; set +a; fi; \
	done; \
	E2E_DISCORD_LIVE=1 go test -count=1 -timeout 120s ./e2etest/ -run TestE2E_DiscordNotifications

# Run Discord notification browser e2e against real webhook URLs from app .envs files.
e2e-browser-discord-live:
	@set -euo pipefail; \
	for f in apps/dashboard/.envs apps/api/.envs apps/billing/.envs; do \
		if [ -f "$$f" ]; then set -a; . "$$f"; set +a; fi; \
	done; \
	E2E_DISCORD_LIVE=1 bash -lc "cd e2etest/browser && npx playwright test -g 'discord notifications fire'"

# Run browser snapshot tests (saves screenshots locally).
e2e-browser-snapshots:
	cd e2etest/browser && npx playwright test -g "snapshots"

# Run browser snapshot tests and upload to Discord.
e2e-browser-snapshots-live:
	@set -euo pipefail; \
	for f in apps/dashboard/.envs apps/api/.envs apps/billing/.envs; do \
		if [ -f "$$f" ]; then set -a; . "$$f"; set +a; fi; \
	done; \
	bash -lc "cd e2etest/browser && npx playwright test -g 'snapshots'"
