# monitoring-agent-client

Reimplementation of https://github.com/infraweavers/monitoring-agent-scripts/blob/main/check_script_via_monitoring-agent.pl in go, to aid scaling and performance from the OMD/naemon instance as the perl implementation will take approximately 300ms for all checks (in our environment) whereas the go implementation takes about 100ms. 

For comparison, nrpe takes approximately 130ms.
