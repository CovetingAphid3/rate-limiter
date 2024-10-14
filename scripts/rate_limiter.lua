local tokens = redis.call("GET", KEYS[1])
if tokens == false then
    tokens = ARGV[2] -- Set tokens to burst limit if not found
else
    tokens = tonumber(tokens)
end

-- Add tokens based on elapsed time
local lastTime = redis.call("GET", KEYS[2])
local currentTime = tonumber(ARGV[1])

if lastTime == false then
    lastTime = currentTime
else
    lastTime = tonumber(lastTime)
end

local elapsed = (currentTime - lastTime) / 1000.0
local tokensToAdd = math.floor(elapsed / ARGV[3])

tokens = math.min(tokens + tokensToAdd, tonumber(ARGV[2]))

-- Save the new timestamp and check if we can allow the request
if tokens > 0 then
    tokens = tokens - 1
    redis.call("SET", KEYS[1], tokens)
    redis.call("SET", KEYS[2], currentTime)
    return 1
else
    redis.call("SET", KEYS[2], currentTime)
    return 0
end

