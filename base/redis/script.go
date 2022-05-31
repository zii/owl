package redis

import "github.com/gomodule/redigo/redis"

// Only increase the value if the key exists.
var lua_incrby = redis.NewScript(1, `
	local key = KEYS[1];
	local amount = tonumber(ARGV[1]);
	if amount < 0 then
		local n = redis.call("GET", key);
		n = tonumber(n or 0);
		if n < -amount then
			return nil;
		else
			n = redis.call("INCRBY", key, amount);
			return n;
		end
	elseif redis.call("EXISTS", key) == 1 then
		local n = redis.call("INCRBY", key, amount);
		return n;
	end
	return nil;
`)
