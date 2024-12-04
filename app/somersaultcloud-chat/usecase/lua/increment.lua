-- increment.lua
local current = redis.call("GET", KEYS[1])
if not current then
    current = 0
end
current = current + 1
redis.call("SET", KEYS[1], current)
return current
