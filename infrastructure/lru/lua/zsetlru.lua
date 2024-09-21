redis.replicate_commands()

local zset_key = KEYS[1]
local member = ARGV[1]
local max_items = tonumber(ARGV[2])

-- 更新 ZSet 信息，使用当前时间戳作为分数
local timestamp = redis.call("TIME")
local score = timestamp[1] + timestamp[2] / 1000000
redis.call("ZADD", zset_key, score, member)

local removed_member = nil

-- 检查 ZSet 大小并移除最旧的元素
local zset_size = redis.call("ZCARD", zset_key)
if zset_size > max_items then
    local range_result = redis.call("ZRANGE", zset_key, 0, 0)
    if range_result and #range_result > 0 then
        removed_member = range_result[1]
        redis.call("ZREM", zset_key, removed_member)
    end
end

if removed_member then
    return removed_member
else
    return nil
end