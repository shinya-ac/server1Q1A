FROM mysql:8.0
ENV LANG ja_JP.UTF-8

# WORKDIR /mysql

COPY ./init/01-create_tasks_table.sh /docker-entrypoint-initdb.d/01-create_tasks_table.sh
COPY ./init/02-create_folders_table.sh /docker-entrypoint-initdb.d/02-create_folders_table.sh