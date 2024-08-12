local list_key = KEYS[1]
local lru_key = KEYS[2]
local value = ARGV[1]
local max_items = tonumber(ARGV[2])

-- 将新元素插入列表左端
redis.call("LPUSH", list_key, value)

-- 更新 LRU 信息，使用当前时间戳作为分数
local timestamp = redis.call("TIME")
local score = timestamp[1] + timestamp[2] / 1000000
redis.call("ZADD", lru_key, score, value)

-- 检查列表大小并移除最旧的元素
local list_size = redis.call("LLEN", list_key)
if list_size > max_items then
    local oldest = redis.call("LRANGE", list_key, -1, -1)[1]
    redis.call("RPOP", list_key)
    redis.call("ZREM", lru_key, oldest)
end

return "OK"