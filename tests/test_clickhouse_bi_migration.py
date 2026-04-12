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

    def test_profile_view_contract_section(self) -> None:
        sql = MIGRATION.read_text(encoding="utf-8").lower()
        profile_section = sql.split(
            "create view if not exists v_bi_subscriber_object_profile as", 1
        )[1].split("-- +goose down", 1)[0]
        outer_select = profile_section.split("from\n(", 1)[0]

        self.assertIn("group by subscriber_id, object_id", profile_section)
        self.assertIn("argmax(subscriber_account_number, finished_at)", profile_section)
        self.assertIn("argmax(subscriber_status, finished_at)", profile_section)
        self.assertIn("argmax(object_address, finished_at)", profile_section)
        self.assertIn("argmax(object_have_automaton, finished_at)", profile_section)
        self.assertNotIn("group by subscriber_id, subscriber_account_number", profile_section)
        self.assertNotIn("group by subscriber_id, object_address", profile_section)
        self.assertNotIn("group by subscriber_id, object_have_automaton", profile_section)
        self.assertIn("object_id,", outer_select)
        for alias in (
            "subscriber_account_number",
            "subscriber_status_ru",
            "object_id",
            "object_address",
            "object_have_automaton",
            "automaton_state_ru",
            "last_task_day",
            "total_tasks_count",
            "violations_detected_count",
            "unauthorized_consumers_count",
        ):
            self.assertIn(alias, outer_select)
        self.assertIn("select\n    subscriber_id,", outer_select)
        self.assertIn("if(object_have_automaton, 'есть автомат', 'нет автомата') as automaton_state_ru", outer_select)
        self.assertIn("last_task_day,", outer_select)
        self.assertIn("total_tasks_count,", outer_select)
        self.assertIn("violations_detected_count,", outer_select)
        self.assertIn("unauthorized_consumers_count", outer_select)
        for forbidden in (
            "subscriber_phone_number",
            "subscriber_email",
            "subscriber_inn",
            "subscriber_birth_date",
        ):
            self.assertNotIn(forbidden, profile_section)

    def test_inspection_results_view_contract_section(self) -> None:
        sql = MIGRATION.read_text(encoding="utf-8").lower()
        inspection_section = sql.split(
            "create view if not exists v_bi_inspection_results as", 1
        )[1].split("create view if not exists v_bi_subscriber_object_profile as", 1)[0]

        self.assertIn("subscriber_status_ru", inspection_section)
        self.assertIn("multiif(", inspection_section)
        self.assertIn("subscriber_status = 'active', 'активен'", inspection_section)
        self.assertIn("subscriber_status = 'violator', 'нарушитель'", inspection_section)
        self.assertIn("subscriber_status = 'archived', 'архивный'", inspection_section)
        self.assertIn("'неизвестно'", inspection_section)
        self.assertIn(
            "group by day, inspection_type_ru, inspection_result_ru, subscriber_status_ru",
            inspection_section,
        )


if __name__ == "__main__":
    unittest.main()
