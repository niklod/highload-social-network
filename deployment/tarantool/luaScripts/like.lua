function find_user_by_name(lastname_pattern, firstname_pattern)
    local result = {}
    local i = 0

    for _, addr in box.space.user.index.lastname_firstname_idx:pairs({lastname_pattern, firstname_pattern}, { iterator = 'GE' }) do
        if string.startswith(addr[3], lastname_pattern, 1, -1) and string.startswith(addr[2], firstname_pattern, 1, -1) then
            table.insert(result, addr)
            i = i + 1
        end

        if i == 100 then
            return result
        end
    end

    return result
end