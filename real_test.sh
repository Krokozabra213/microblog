#!/bin/bash

echo "========================================="
echo "🚀 Полный тест Microblog API"
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

# Функция для извлечения значения поля из JSON
extract_field() {
    local json="$1"
    local field="$2"
    echo "$json" | grep -o "\"$field\":\"[^\"]*\"" | head -1 | sed "s/\"$field\":\"\([^\"]*\)\"/\1/"
}

# 1. Создаём пользователя Alice
echo "1️⃣  Создаём пользователя Alice..."
USER1_RESPONSE=$(curl -s -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d '{"username":"alice"}')
echo "$USER1_RESPONSE"
USER1_ID=$(extract_field "$USER1_RESPONSE" "id")
echo "   Alice ID: $USER1_ID"
echo ""

sleep 0.5

# 2. Создаём пользователя Bob
echo "2️⃣  Создаём пользователя Bob..."
USER2_RESPONSE=$(curl -s -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d '{"username":"bob"}')
echo "$USER2_RESPONSE"
USER2_ID=$(extract_field "$USER2_RESPONSE" "id")
echo "   Bob ID: $USER2_ID"
echo ""

sleep 0.5

# 3. Alice создаёт пост
echo "3️⃣  Alice создаёт пост..."
POST1_RESPONSE=$(curl -s -X POST "$BASE_URL/posts" \
  -H "Content-Type: application/json" \
  -d "{\"author_id\":\"$USER1_ID\",\"text\":\"Привет, это мой первый пост!\"}")
echo "$POST1_RESPONSE"
POST1_ID=$(extract_field "$POST1_RESPONSE" "id")
echo "   Post ID: $POST1_ID"
echo ""

sleep 0.5

# 4. Bob создаёт пост
echo "4️⃣  Bob создаёт пост..."
POST2_RESPONSE=$(curl -s -X POST "$BASE_URL/posts" \
  -H "Content-Type: application/json" \
  -d "{\"author_id\":\"$USER2_ID\",\"text\":\"Контент от Bob!\"}")
echo "$POST2_RESPONSE"
POST2_ID=$(extract_field "$POST2_RESPONSE" "id")
echo "   Post ID: $POST2_ID"
echo ""

sleep 0.5

# 5. Получаем все посты
echo "5️⃣  Получаем все посты..."
ALL_POSTS=$(curl -s -X GET "$BASE_URL/posts")
echo "$ALL_POSTS"
echo ""

sleep 0.5

# 6. Получаем конкретный пост Alice
echo "6️⃣  Получаем пост Alice по ID ($POST1_ID)..."
ALICE_POST=$(curl -s -X GET "$BASE_URL/posts/$POST1_ID")
echo "$ALICE_POST"
echo ""

sleep 0.5

# 7. Bob лайкает пост Alice
echo "7️⃣  Bob лайкает пост Alice..."
LIKE1_RESPONSE=$(curl -s -X POST "$BASE_URL/posts/$POST1_ID/like" \
  -H "Content-Type: application/json" \
  -d "{\"user_id\":\"$USER2_ID\"}")
echo "$LIKE1_RESPONSE"
echo ""

sleep 0.5

# 8. Alice лайкает свой пост
echo "8️⃣  Alice лайкает свой пост..."
LIKE2_RESPONSE=$(curl -s -X POST "$BASE_URL/posts/$POST1_ID/like" \
  -H "Content-Type: application/json" \
  -d "{\"user_id\":\"$USER1_ID\"}")
echo "$LIKE2_RESPONSE"
echo ""

sleep 0.5

# 9. Проверяем финальное состояние поста Alice
echo "9️⃣  Финальное состояние поста Alice (должно быть 2 лайка)..."
FINAL_POST=$(curl -s -X GET "$BASE_URL/posts/$POST1_ID")
echo "$FINAL_POST"
# Подсчитываем количество лайков (количество вхождений в массив likes)
LIKES_COUNT=$(echo "$FINAL_POST" | grep -o "\"likes\":\[" | wc -l)
if echo "$FINAL_POST" | grep -q "\"likes\":\[\""; then
    LIKES_COUNT=$(echo "$FINAL_POST" | grep -o "\"[a-f0-9-]\{36\}\"" | grep -A 10 "\"likes\":\[" | wc -l)
fi
echo "   Лайков: проверьте визуально (в массиве likes)"
echo ""

sleep 0.5

# 10. Попытка повторного лайка (должна быть ошибка)
echo "🔟 Bob пытается лайкнуть второй раз (должна быть ошибка)..."
DUPLICATE_LIKE=$(curl -s -X POST "$BASE_URL/posts/$POST1_ID/like" \
  -H "Content-Type: application/json" \
  -d "{\"user_id\":\"$USER2_ID\"}")
echo "$DUPLICATE_LIKE"
echo ""

sleep 0.5

# 11. Попытка создать пост с пустым текстом
echo "1️⃣1️⃣ Попытка создать пост с пустым текстом (должна быть ошибка)..."
EMPTY_POST=$(curl -s -X POST "$BASE_URL/posts" \
  -H "Content-Type: application/json" \
  -d "{\"author_id\":\"$USER1_ID\",\"text\":\"\"}")
echo "$EMPTY_POST"
echo ""

sleep 0.5

# 12. Попытка лайкнуть несуществующий пост
echo "1️⃣2️⃣ Попытка лайкнуть несуществующий пост (должна быть ошибка)..."
FAKE_LIKE=$(curl -s -X POST "$BASE_URL/posts/fake-post-id/like" \
  -H "Content-Type: application/json" \
  -d "{\"user_id\":\"$USER1_ID\"}")
echo "$FAKE_LIKE"
echo ""

# 13. Получение поста Bob (должно быть 0 лайков)
echo "1️⃣3️⃣ Получение поста Bob (должно быть 0 лайков)..."
BOB_POST=$(curl -s -X GET "$BASE_URL/posts/$POST2_ID")
echo "$BOB_POST"
echo ""

echo "========================================="
echo "✅ Тестирование завершено!"
echo "========================================="
echo ""
echo "📊 Сводка:"
echo "   Пользователи: Alice ($USER1_ID), Bob ($USER2_ID)"
echo "   Посты: Alice ($POST1_ID), Bob ($POST2_ID)"
echo ""
echo "🎯 Проверьте результаты выше:"
echo "   - Пост Alice должен иметь 2 лайка (в массиве likes)"
echo "   - Пост Bob должен иметь 0 лайков (пустой массив)"
echo "   - Повторный лайк должен вернуть ошибку"
echo "   - Пустой текст поста должен вернуть ошибку"
echo ""
