function add_deprecated_headers(deprecated_at, deprecated_link, deprecated_sunset)
    -- Verifica se os 3 parâmetros estão presentes
    if deprecated_at and deprecated_link and deprecated_sunset then
        local r = response.load()

        -- Adiciona os headers com as informações de depreciação
        r:headers("x-deprecated-at", deprecated_at)
        r:headers("x-deprecated-link", deprecated_link)
        r:headers("x-deprecated-sunset", deprecated_sunset)

        print("[Lua Plugin] Headers de depreciação adicionados com sucesso.")
    else
        -- Se algum parâmetro estiver ausente, printa um log e não adiciona os headers
        print("[Lua Plugin] Informações de depreciação incompletas. Headers não adicionados.")
    end
end
