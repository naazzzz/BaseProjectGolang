# ──────────────────────────────────────────────────────────────
# Тесты без docker-compose — только docker build + run
# ──────────────────────────────────────────────────────────────

# Имя образа и тега (можно переопределить через переменные)
TEST_IMAGE      ?= listing-test
LINT_IMAGE   ?= listing-lint
TEST_CONTAINER  ?= listing-test-run

# 1. Сборка тестового образа (из DockerfileForTest)
test-build:
	docker build -f DockerfileForTest -t $(TEST_IMAGE):latest .

# 2. Запуск тестов (обычных) — просто прогон
test: test-build
	docker run --rm $(TEST_IMAGE):latest

# 3. Запуск тестов с покрытием и артефактами (вот это нужно в CI)
test-coverage: test-build
	@echo "Запуск тестов с генерацией отчётов (junit + cobertura)..."
	@mkdir -p coverage
	docker run --rm \
		-v $(PWD)/coverage:/listing/coverage \
		$(TEST_IMAGE):latest

# 4. Быстрая команда: пересобрать + запустить с покрытием + почистить контейнер
test-quick: test-build test-coverage
	@echo "Тесты завершены. Отчёты лежат в ./coverage/"

# 5. Очистка образа
test-clean:
	docker rmi $(TEST_IMAGE):latest || true
	docker volume prune -f || true

# 6. Полная перезагрузка
test-reset: test-clean test-quick

# ──────────────────────────────────────────────────────────────
# ЛИНТ (через DockerfileForLint)
# ──────────────────────────────────────────────────────────────

# Сборка образа с golangci-lint
lint-build:
	docker build -f DockerfileForLint -t $(LINT_IMAGE):latest .

# Запуск линта (падает при ошибках — как и должно быть)
lint: lint-build
	@echo "Запуск golangci-lint..."
	docker run --rm \
		-v $(PWD):/listing \
		-w /listing \
		$(LINT_IMAGE):latest

# Линт с автоматическим исправлением (если поддерживается линтерами)
lint-fix: lint-build
	@echo "Запуск golangci-lint с автофиксом..."
	docker run --rm \
		-v $(PWD):/listing \
		-w /listing \
		$(LINT_IMAGE):latest golangci-lint run ./... --fix

# Очистка линт-образа
lint-clean:
	docker rmi $(LINT_IMAGE):latest || true

# ──────────────────────────────────────────────────────────────
# Комбинированные и удобные цели
# ──────────────────────────────────────────────────────────────

# Всё по-чистому: линт → тесты с покрытием
ci: lint test-coverage

# Полная проверка перед коммитом/мерджем
check: lint test-quick
	@echo "Всё ок! Линт прошёл, тесты прошли, покрытие сгенерировано."

# Полная очистка всех образов
clean: test-clean lint-clean
	docker image prune -f

# Пересборка всего с нуля
reset: clean ci
