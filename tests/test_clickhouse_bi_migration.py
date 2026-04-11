from pathlib import Path
import unittest


ROOT = Path(__file__).resolve().parents[1]
MIGRATION = ROOT / "database" / "migrations" / "clickhouse" / "00002_bi_views.sql"


class ClickHouseBiMigrationContractTests(unittest.TestCase):
    def test_migration_file_exists(self) -> None:
        self.assertTrue(MIGRATION.exists(), f"missing migration: {MIGRATION}")

    def test_migration_declares_all_views(self) -> None:
        sql = MIGRATION.read_text(encoding="utf-8").lower()
        for view_name in (
            "v_bi_tasks_daily",
            "v_bi_brigade_performance",
            "v_bi_inspection_results",
            "v_bi_subscriber_object_profile",
        ):
            self.assertIn(f"create view if not exists {view_name}", sql)

    def test_profile_view_uses_stable_keys_and_latest_fields(self) -> None:
        sql = MIGRATION.read_text(encoding="utf-8").lower()
        profile_section = sql.split(
            "create view if not exists v_bi_subscriber_object_profile as", 1
        )[1].split("-- +goose down", 1)[0]

        self.assertIn("group by subscriber_id, object_id", profile_section)
        self.assertIn("argmax(subscriber_account_number, finished_at)", profile_section)
        self.assertIn("argmax(subscriber_status, finished_at)", profile_section)
        self.assertIn("argmax(object_address, finished_at)", profile_section)
        self.assertIn("argmax(object_have_automaton, finished_at)", profile_section)
        self.assertNotIn("group by subscriber_id, subscriber_account_number", profile_section)
        self.assertNotIn("group by subscriber_id, object_address", profile_section)
        self.assertNotIn("group by subscriber_id, object_have_automaton", profile_section)

    def test_profile_view_does_not_expose_pii_columns(self) -> None:
        sql = MIGRATION.read_text(encoding="utf-8").lower()
        for forbidden in (
            "subscriber_phone_number",
            "subscriber_email",
            "subscriber_inn",
            "subscriber_birth_date",
        ):
            self.assertNotIn(forbidden, sql)


if __name__ == "__main__":
    unittest.main()
