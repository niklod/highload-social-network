function batch_insert(list)
    box.begin()
        for _, record in ipairs(list) do
            box.space.user:replace(record)
        end
    box.commit()
end