local key = KEYS[1]
local cntKey = key .. ":cnt"
local val = ARGS[1]

local ttl = tonumber(redis.call("ttl", key))

if ttl == -1 then
    -- key exists but does not have expire time
    return -2
elseif ttl == -2 or ttl < 540 then
    -- -2 means key does not exists and less than 9 min, then send code
    redis.call("set", key, val)
    redis.call("expire", key, 600)
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, 600)
else
    -- code sent too much times within 10 min
    return -1
end