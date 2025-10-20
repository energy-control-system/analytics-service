select id, report_id, file_id, created_at
from attachments
where report_id = any ($1);
