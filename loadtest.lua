transact = function()
    headers = {}
    headers["Content-Type"] = "application/json"
    id = math.random(1, 5)
    value = math.random(1, 10000)
    tp = "cd"
    tpidx = math.random(1,2)
    body = '{"valor":' .. tostring(value) .. ',"tipo":"' .. string.sub(tp, tpidx, tpidx) .. '","descricao":"crebito"}'
    return wrk.format("POST", "/clientes/" .. tostring(id) .. "/transacoes", headers, body)
end

history = function()
    headers = {}
    headers["Accept-Encoding"] = "gzip"
    id = math.random(1, 6)
    body = ""
    return wrk.format("GET", "/clients/" .. tostring(id) .. "/extrato", headers, body)
end


requests = {}
requests[0] = transact
requests[1] = history
requests[2] = transact
requests[3] = history
requests[4] = history
requests[5] = history
requests[6] = history

request = function()
    return requests[math.random(0,6)]()
end
