import datetime

from airflow import DAG
from airflow.operators.empty import EmptyOperator
from airflow.operators.bash import BashOperator


scrap_dag = DAG(
    dag_id="fresher_scout",
    start_date=datetime.datetime(2025, 3, 3),
    schedule="@daily",
    catchup=False,
)

start = EmptyOperator(taks_id="scrap", dag=scrap_dag)

run_scout = BashOperator(
    task_id="run_scout",
    bash_command="./bin/scout",
    dag=scrap_dag,
)

start >> run_scout
