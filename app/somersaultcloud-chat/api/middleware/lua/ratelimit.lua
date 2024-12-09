local key = KEYS[1]
local rate = tonumber(ARGV[1])
local capacity = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local requested = tonumber(ARGV[4])

local tokens = tonumber(redis.call("HGET", key, "tokens") or capacity)
local last_time = tonumber(redis.call("HGET", key, "last_time") or now)

local delta = math.max((now - last_time) * rate / 1000, 0)
tokens = math.min(tokens + delta, capacity)

if tokens >= requested then
    tokens = tokens - requested
    redis.call("HSET", key, "tokens", tokens)
    redis.call("HSET", key, "last_time", now)
    return 1
else
    redis.call("HSET", key, "tokens", tokens)
    redis.call("HSET", key, "last_time", now)
    return 0
end