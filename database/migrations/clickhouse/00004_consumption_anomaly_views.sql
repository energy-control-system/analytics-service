-- +goose Up
alter table finished_tasks
    add column if not exists inspected_devices Array(Tuple(
        id Int64,
        device_id Int64,
        value Decimal(15, 2),
        consumption_kwh Decimal(15, 2),
        created_at DateTime('UTC')
    )) default [];

drop view if exists v_bi_consumption_anomalies;
drop view if exists v_bi_consumption_monthly;

create view if not exists v_bi_consumption_monthly as
select
    toStartOfMonth(finished_at) as month,
    subscriber_id,
    subscriber_account_number,
    concat(subscriber_surname, ' ', subscriber_name, ' ', subscriber_patronymic) as subscriber_full_name,
    object_id,
    object_address,
    replaceRegexpOne(object_address, ',.*$', '') as district_name,
    groupUniqArray(toString(device_reading.1)) as inspected_device_ids,
    groupUniqArray(toString(device_reading.2)) as device_ids,
    sum(toDecimal64(device_reading.4, 2)) as monthly_consumption_kwh,
    count() as readings_count,
    max(finished_at) as last_reading_at
from finished_tasks
array join inspected_devices as device_reading
where toDecimal64(device_reading.4, 2) > 0
group by
    month,
    subscriber_id,
    subscriber_account_number,
    subscriber_full_name,
    object_id,
    object_address,
    district_name;

create view if not exists v_bi_consumption_anomalies as
with scored as
(
    select
        month,
        subscriber_id,
        subscriber_account_number,
        subscriber_full_name,
        object_id,
        object_address,
        district_name,
        device_ids,
        monthly_consumption_kwh,
        readings_count,
        last_reading_at,
        subscriber_avg_consumption_kwh,
        subscriber_months_count,
        district_avg_consumption_kwh,
        if(
            subscriber_avg_consumption_kwh > 0,
            abs(monthly_consumption_kwh - subscriber_avg_consumption_kwh) / subscriber_avg_consumption_kwh,
            0
        ) as subscriber_deviation_ratio,
        if(
            district_avg_consumption_kwh > 0,
            monthly_consumption_kwh / district_avg_consumption_kwh,
            0
        ) as district_deviation_ratio
    from
    (
        select
            month,
            subscriber_id,
            subscriber_account_number,
            subscriber_full_name,
            object_id,
            object_address,
            district_name,
            device_ids,
            toFloat64(monthly_consumption_kwh) as monthly_consumption_kwh,
            readings_count,
            last_reading_at,
            avg(toFloat64(monthly_consumption_kwh)) over (partition by subscriber_id, object_id) as subscriber_avg_consumption_kwh,
            count() over (partition by subscriber_id, object_id) as subscriber_months_count,
            (
                sum(toFloat64(monthly_consumption_kwh)) over (partition by district_name)
                - toFloat64(monthly_consumption_kwh)
            ) / nullIf(count() over (partition by district_name) - 1, 0) as district_avg_consumption_kwh
        from v_bi_consumption_monthly
    )
)
select
    month,
    subscriber_id,
    subscriber_account_number,
    subscriber_full_name,
    object_id,
    object_address,
    district_name,
    device_ids,
    monthly_consumption_kwh,
    round(subscriber_avg_consumption_kwh, 2) as subscriber_avg_consumption_kwh,
    subscriber_months_count,
    round(district_avg_consumption_kwh, 2) as district_avg_consumption_kwh,
    round(subscriber_deviation_ratio * 100, 2) as subscriber_deviation_percent,
    round((district_deviation_ratio - 1) * 100, 2) as district_deviation_percent,
    multiIf(
        subscriber_months_count >= 3
            and monthly_consumption_kwh > subscriber_avg_consumption_kwh
            and subscriber_deviation_ratio >= 0.5,
        'Скачок относительно истории абонента',
        subscriber_months_count >= 3
            and monthly_consumption_kwh < subscriber_avg_consumption_kwh
            and subscriber_deviation_ratio >= 0.5,
        'Провал относительно истории абонента',
        district_deviation_ratio >= 2.5,
        'Выше среднего по району',
        'Норма'
    ) as anomaly_reason,
    greatest(subscriber_deviation_ratio, district_deviation_ratio - 1) as severity_score,
    readings_count,
    last_reading_at
from scored
where
    (subscriber_months_count >= 3 and subscriber_deviation_ratio >= 0.5)
    or district_deviation_ratio >= 2.5;

-- +goose Down
drop view if exists v_bi_consumption_anomalies;
drop view if exists v_bi_consumption_monthly;

alter table finished_tasks
    drop column if exists inspected_devices;
