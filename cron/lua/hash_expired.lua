local hash_key = KEYS[1]
local expire_key = KEYS[2]
local field = ARGV[1]
local value = ARGV[2]
local ttl = tonumber(ARGV[3])

-- 存储字段和值
redis.call("HSET", hash_key, field, value)

-- 设置过期时间（Unix 时间戳）
local expire_at = redis.call("TIME")[1] + ttl
redis.call("HSET", expire_key, field, expire_at)

return "OK"
