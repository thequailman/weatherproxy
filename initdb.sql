create database weatherproxy;
\c weatherproxy
create role "weatherproxy" with createdb login password 'weatherproxy';
revoke all privileges on schema public from public;
grant all privileges on database weatherproxy to "weatherproxy";
grant all privileges on schema public to "weatherproxy";
