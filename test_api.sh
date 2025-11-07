#!/bin/bash

# Тестовый скрипт для проверки API microblog

BASE_URL="http://localhost:8080"

echo "=================================="
echo "Тестирование API Microblog"
echo "=================================="
echo ""

# Цвета для вывода
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}1. Регистрация первого пользователя (Alice)${NC}"
USER1_RESPONSE=$(curl -s -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d '{"username":"alice"}')
echo "$USER1_RESPONSE" | jq '.'
USER1_ID=$(echo "$USER1_RESPONSE" | jq -r '.id')
echo -e "${GREEN}User1 ID: $USER1_ID${NC}"
echo ""

sleep 1

echo -e "${YELLOW}2. Регистрация второго пользователя (Bob)${NC}"
USER2_RESPONSE=$(curl -s -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d '{"username":"bob"}')
echo "$USER2_RESPONSE" | jq '.'
USER2_ID=$(echo "$USER2_RESPONSE" | jq -r '.id')
echo -e "${GREEN}User2 ID: $USER2_ID${NC}"
echo ""

sleep 1

echo -e "${YELLOW}3. Попытка создать пользователя с пустым username (должна быть ошибка)${NC}"
curl -s -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d '{"username":""}' | jq '.'
echo ""

sleep 1

echo -e "${YELLOW}4. Создание поста от Alice${NC}"
POST1_RESPONSE=$(curl -s -X POST "$BASE_URL/posts" \
  -H "Content-Type: application/json" \
  -d "{\"author_id\":\"$USER1_ID\",\"text\":\"Hello, this is my first post!\"}")
echo "$POST1_RESPONSE" | jq '.'
POST1_ID=$(echo "$POST1_RESPONSE" | jq -r '.id')
echo -e "${GREEN}Post1 ID: $POST1_ID${NC}"
echo ""

sleep 1

echo -e "${YELLOW}5. Создание поста от Bob${NC}"
POST2_RESPONSE=$(curl -s -X POST "$BASE_URL/posts" \
  -H "Content-Type: application/json" \
  -d "{\"author_id\":\"$USER2_ID\",\"text\":\"Bob's amazing content here!\"}")
echo "$POST2_RESPONSE" | jq '.'
POST2_ID=$(echo "$POST2_RESPONSE" | jq -r '.id')
echo -e "${GREEN}Post2 ID: $POST2_ID${NC}"
echo ""

sleep 1

echo -e "${YELLOW}6. Попытка создать пост с пустым текстом (должна быть ошибка)${NC}"
curl -s -X POST "$BASE_URL/posts" \
  -H "Content-Type: application/json" \
  -d "{\"author_id\":\"$USER1_ID\",\"text\":\"\"}" | jq '.'
echo ""

sleep 1

echo -e "${YELLOW}7. Получение всех постов${NC}"
curl -s -X GET "$BASE_URL/posts" | jq '.'
echo ""

sleep 1

echo -e "${YELLOW}8. Получение конкретного поста по ID${NC}"
curl -s -X GET "$BASE_URL/posts/$POST1_ID" | jq '.'
echo ""

sleep 1

echo -e "${YELLOW}9. Bob лайкает пост Alice${NC}"
curl -s -X POST "$BASE_URL/posts/$POST1_ID/like" \
  -H "Content-Type: application/json" \
  -d "{\"user_id\":\"$USER2_ID\"}" | jq '.'
echo ""

sleep 1

echo -e "${YELLOW}10. Alice лайкает пост Alice${NC}"
curl -s -X POST "$BASE_URL/posts/$POST1_ID/like" \
  -H "Content-Type: application/json" \
  -d "{\"user_id\":\"$USER1_ID\"}" | jq '.'
echo ""

sleep 1

echo -e "${YELLOW}11. Попытка Bob лайкнуть пост второй раз (должна быть ошибка)${NC}"
curl -s -X POST "$BASE_URL/posts/$POST1_ID/like" \
  -H "Content-Type: application/json" \
  -d "{\"user_id\":\"$USER2_ID\"}" | jq '.'
echo ""

sleep 1

echo -e "${YELLOW}12. Проверка финального состояния поста с лайками${NC}"
curl -s -X GET "$BASE_URL/posts/$POST1_ID" | jq '.'
echo ""

sleep 1

echo -e "${YELLOW}13. Получение несуществующего поста (должна быть ошибка)${NC}"
curl -s -X GET "$BASE_URL/posts/nonexistent-id" | jq '.'
echo ""

echo "=================================="
echo -e "${GREEN}Тестирование завершено!${NC}"
echo "=================================="
