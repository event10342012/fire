local key = KEYS[1]
local cntKey = key .. ":cnt"

local expectedCode = ARGS[1]

local cnt = tonumber(redis.call("Get", cntKey))
local code = redis.call("Get", key)

if cnt == nil or cnt <= 0 then
    -- run out of the verify times
    return -1
end

if code == expectedCode then
    -- input correct, and set cntKey = 0 to avoid input again
    redis.call("Set", cntKey, 0)
    return 0
else
    redis.call("decr", cntKey)
    -- not equal, input incorrectly
    return -2
end
