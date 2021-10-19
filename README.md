# monitoring-agent-client

This is a go reimplementation of the perl [check_script_via_monitoring-agent.pl](https://github.com/infraweavers/monitoring-agent-scripts/blob/main/check_script_via_monitoring-agent.pl). 

This should have comparable performance/load on the monitoring-server (i.e. Nagios/Naemon/OMD), in our testing a check was ~110ms vs ~130ms for nrpe.

The perl implementation has a significantly higher cost of execution (approx 300ms) so if there is an issue with load on the monitoring server itself, it would be wise to transition to using this implementation.
