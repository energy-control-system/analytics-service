insert into attachments (report_id, file_id)
values (:report_id, :file_id)
returning id, report_id, file_id, created_at;
