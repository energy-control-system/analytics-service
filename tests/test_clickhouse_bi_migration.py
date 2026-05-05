from pathlib import Path
import unittest


ROOT = Path(__file__).resolve().parents[1]
MIGRATION_02 = ROOT / "database" / "migrations" / "clickhouse" / "00002_bi_views.sql"
MIGRATION_03 = ROOT / "database" / "migrations" / "clickhouse" / "00003_bi_inspection_results.sql"
MIGRATION_04 = ROOT / "database" / "migrations" / "clickhouse" / "00004_consumption_anomaly_views.sql"
ADD_FINISHED_TASK_SQL = ROOT / "database" / "analytics" / "sql" / "add_finished_task.sql"


class ClickHouseBiMigrationContractTests(unittest.TestCase):
    def test_migration_file_exists(self) -> None:
        self.assertTrue(MIGRATION_02.exists(), f"missing migration: {MIGRATION_02}")

    def test_upgrade_migration_file_exists(self) -> None:
        self.assertTrue(MIGRATION_03.exists(), f"missing migration: {MIGRATION_03}")

    def test_consumption_anomaly_migration_file_exists(self) -> None:
        self.assertTrue(MIGRATION_04.exists(), f"missing migration: {MIGRATION_04}")

    def test_migration_declares_all_views(self) -> None:
        sql = MIGRATION_02.read_text(encoding="utf-8").lower()
        for view_name in (
            "v_bi_tasks_daily",
            "v_bi_brigade_performance",
            "v_bi_inspection_results",
            "v_bi_subscriber_object_profile",
        ):
            self.assertIn(f"create view if not exists {view_name}", sql)

    def test_profile_view_contract_section(self) -> None:
        sql = MIGRATION_02.read_text(encoding="utf-8").lower()
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
        sql = MIGRATION_02.read_text(encoding="utf-8").lower()
        inspection_section = sql.split(
            "create view if not exists v_bi_inspection_results as", 1
        )[1].split("create view if not exists v_bi_subscriber_object_profile as", 1)[0]

        self.assertIn(
            "group by day, inspection_type_ru, inspection_result_ru",
            inspection_section,
        )
        self.assertNotIn("subscriber_status_ru", inspection_section)
        self.assertNotIn("day_tasks_share_ratio", inspection_section)
        self.assertNotIn("sum(tasks_count) over (partition by day)", inspection_section)

    def test_upgrade_inspection_results_view_contract_section(self) -> None:
        sql = MIGRATION_03.read_text(encoding="utf-8").lower()
        self.assertIn("drop view if exists v_bi_inspection_results", sql)
        inspection_section = sql.split(
            "create view if not exists v_bi_inspection_results as", 1
        )[1].split("-- +goose down", 1)[0]

        self.assertIn("subscriber_status_ru", inspection_section)
        self.assertIn("day_tasks_share_ratio", inspection_section)
        self.assertIn("multiif(", inspection_section)
        self.assertIn("subscriber_status = 'active', 'активен'", inspection_section)
        self.assertIn("subscriber_status = 'violator', 'нарушитель'", inspection_section)
        self.assertIn("subscriber_status = 'archived', 'архивный'", inspection_section)
        self.assertIn("'неизвестно'", inspection_section)
        self.assertIn("sum(tasks_count) over (partition by day)", inspection_section)
        self.assertIn(
            "group by day, inspection_type_ru, inspection_result_ru, subscriber_status_ru",
            inspection_section,
        )

    def test_upgrade_inspection_results_rollback_contract(self) -> None:
        sql = MIGRATION_03.read_text(encoding="utf-8").lower()
        down_section = sql.split("-- +goose down", 1)[1]

        self.assertIn("drop view if exists v_bi_inspection_results", down_section)
        self.assertIn("create view if not exists v_bi_inspection_results", down_section)
        self.assertIn(
            "group by day, inspection_type_ru, inspection_result_ru",
            down_section,
        )
        self.assertNotIn("subscriber_status_ru", down_section)
        self.assertNotIn("day_tasks_share_ratio", down_section)
        self.assertNotIn("sum(tasks_count) over (partition by day)", down_section)

    def test_consumption_anomaly_views_contract(self) -> None:
        sql = MIGRATION_04.read_text(encoding="utf-8").lower()

        self.assertIn("alter table finished_tasks", sql)
        self.assertIn("add column if not exists inspected_devices", sql)
        self.assertIn("array(tuple(", sql)
        self.assertIn("consumption_kwh decimal(15, 2)", sql)
        self.assertIn("create view if not exists v_bi_consumption_monthly", sql)
        self.assertIn("create view if not exists v_bi_consumption_anomalies", sql)
        self.assertIn("array join inspected_devices", sql)
        self.assertIn("tostring(device_reading.1)", sql)
        self.assertIn("sum(todecimal64(device_reading.4, 2))", sql)
        self.assertIn("avg(tofloat64(monthly_consumption_kwh)) over (partition by subscriber_id, object_id)", sql)
        self.assertIn("count() over (partition by subscriber_id, object_id)", sql)
        self.assertIn("partition by district_name", sql)
        self.assertIn("subscriber_months_count >= 3", sql)
        self.assertIn("subscriber_deviation_ratio >= 0.5", sql)
        self.assertIn("district_deviation_ratio >= 2.5", sql)
        self.assertIn("district_deviation_ratio <= 0.4", sql)
        self.assertIn("'скачок относительно истории абонента'", sql)
        self.assertIn("'провал относительно истории абонента'", sql)
        self.assertIn("'выше среднего по району'", sql)
        self.assertIn("'ниже среднего по району'", sql)

    def test_finished_task_insert_names_inspected_devices_column(self) -> None:
        sql = ADD_FINISHED_TASK_SQL.read_text(encoding="utf-8").lower()

        self.assertIn("insert into finished_tasks", sql)
        self.assertIn("inspection_energy_action_at", sql)
        self.assertIn("inspected_devices", sql)
        self.assertIn("brigade_id", sql)
        self.assertLess(sql.index("inspection_energy_action_at"), sql.index("inspected_devices"))
        self.assertLess(sql.index("inspected_devices"), sql.index("brigade_id"))


if __name__ == "__main__":
    unittest.main()
