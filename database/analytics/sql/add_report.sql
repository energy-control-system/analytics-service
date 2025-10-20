insert into reports (type, period_start, period_end)
values (:type, :period_start, :period_end)
returning id, type, period_start, period_end, created_at;
