[program:demo-supervisord]
directory=/Users/jasper/projects/GolandProjects/go-be-professor/bestpractice/monitoring/demo-supervisord
command=go run demo-supervisord 2>/Users/jasper/projects/GolandProjects/go-be-professor/bestpractice/monitoring/demo-supervisord/gc.log
stdout_logfile=/Users/jasper/projects/GolandProjects/go-be-professor/bestpractice/monitoring/demo-supervisord/stdout.log
stdout_logfile_backups=50
redirect_stderr=true
autostart=true
autorestart=true