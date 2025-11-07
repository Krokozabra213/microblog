#!/bin/bash

echo "========================================="
echo "🚀 Быстрый тест Microblog API"
echo "========================================="
echo ""

BASE_URL="http://localhost:8080"

# Проверяем, запущен ли сервер
echo "📡 Проверяем доступность сервера..."
if ! curl -s -f "$BASE_URL/posts" > /dev/null 2>&1; then
    echo "❌ Сервер не запущен на $BASE_URL"
    echo "Запустите: go run cmd/main.go"
    exit 1
fi
echo "✅ Сервер доступен"
echo ""

# 1. Создаём пользователя Alice
echo "1️⃣  Создаём пользователя Alice..."
USER1=$(curl -s -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d '{"username":"alice"}')
echo "$USER1"
USER1_ID=$(echo "$USER1" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
echo "   ID: $USER1_ID"
echo ""

# 2. Создаём пользователя Bob
echo "2️⃣  Создаём пользователя Bob..."
USER2=$(curl -s -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d '{"username":"bob"}')
echo "$USER2"
USER2_ID=$(echo "$USER2" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
echo "   ID: $USER2_ID"
echo ""

# 3. Alice создаёт пост
echo "3️⃣  Alice создаёт пост..."
POST1=$(curl -s -X POST "$BASE_URL/posts" \
  -H "Content-Type: application/json" \
  -d "{\"author_id\":\"$USER1_ID\",\"text\":\"Привет, это мой первый пост!\"}")
echo "$POST1"
POST1_ID=$(echo "$POST1" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
echo "   ID поста: $POST1_ID"
echo ""

# 4. Bob создаёт пост
echo "4️⃣  Bob создаёт пост..."
POST2=$(curl -s -X POST "$BASE_URL/posts" \
  -H "Content-Type: application/json" \
  -d "{\"author_id\":\"$USER2_ID\",\"text\":\"Bob's контент здесь!\"}")
echo "$POST2"
echo ""

# 5. Получаем все посты
echo "5️⃣  Получаем все посты..."
curl -s -X GET "$BASE_URL/posts"
echo ""
echo ""

# 6. Получаем конкретный пост
echo "6️⃣  Получаем пост Alice по ID..."
curl -s -X GET "$BASE_URL/posts/$POST1_ID"
echo ""
echo ""

# 7. Bob лайкает пост Alice
echo "7️⃣  Bob лайкает пост Alice..."
curl -s -X POST "$BASE_URL/posts/$POST1_ID/like" \
  -H "Content-Type: application/json" \
  -d "{\"user_id\":\"$USER2_ID\"}"
echo ""
echo ""

# 8. Alice лайкает свой пост
echo "8️⃣  Alice лайкает свой пост..."
curl -s -X POST "$BASE_URL/posts/$POST1_ID/like" \
  -H "Content-Type: application/json" \
  -d "{\"user_id\":\"$USER1_ID\"}"
echo ""
echo ""

# 9. Проверяем финальное состояние (должно быть 2 лайка)
echo "9️⃣  Финальное состояние поста Alice (должно быть 2 лайка)..."
curl -s -X GET "$BASE_URL/posts/$POST1_ID"
echo ""
echo ""

# 10. Попытка повторного лайка (должна быть ошибка)
echo "🔟 Bob пытается лайкнуть второй раз (должна быть ошибка)..."
curl -s -X POST "$BASE_URL/posts/$POST1_ID/like" \
  -H "Content-Type: application/json" \
  -d "{\"user_id\":\"$USER2_ID\"}"
echo ""
echo ""

echo "========================================="
echo "✅ Тестирование завершено!"
echo "========================================="
