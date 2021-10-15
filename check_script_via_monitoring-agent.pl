#!/usr/bin/perl
use strict;
use JSON::XS;
use Monitoring::Plugin;
use HTTP::Tiny;

my $plugin = Monitoring::Plugin->new (
	usage => '',
	plugin => $0,
	shortname => 'Check via monitoring-agent',
	blurb => 'Checks via monitoring-agent',
	timeout => "10s",
);

$plugin->add_arg(spec => 'template|t=s', help => 'pnp4nagios template', required => 0);
$plugin->add_arg(spec => 'hostname|h=s', help => 'hostname or ip', required => 1);
$plugin->add_arg(spec => 'port|p=i', help => 'port number', required => 1);
$plugin->add_arg(spec => 'cacert|e=s', help => 'CA certificate', required => 0);
$plugin->add_arg(spec => 'certificate|c=s', help => 'certificate file', required => 0);
$plugin->add_arg(spec => 'key|k=s', help => 'key file', required => 0);
$plugin->add_arg(spec => 'username|u=s', help => 'username', required => 0);
$plugin->add_arg(spec => 'password|p=s', help => 'password', required => 0);
$plugin->add_arg(spec => 'executable|e=s', help => 'executable path', required => 1);
$plugin->add_arg(spec => 'executableArg=s@', help => 'executable arg for multiple specify multiple times', required => 0);
$plugin->add_arg(spec => 'script|s=s', help => 'script location', required => 1);
$plugin->add_arg(spec => 'timeout|t=s', help => 'timeout (e.g. 10s)', required => 0);

$plugin->getopts;

my $username = exists($ENV{'MONITORING_AGENT_USERNAME'}) ? $ENV{'MONITORING_AGENT_USERNAME'} : $plugin->opts->username;
my $password = exists($ENV{'MONITORING_AGENT_PASSWORD'}) ? $ENV{'MONITORING_AGENT_PASSWORD'} : $plugin->opts->password;

my $http = HTTP::Tiny->new("keep_alive" => "0", "verify_SSL" => 1, "SSL_options" => {
	SSL_ca_file => $plugin->opts->cacert,
	SSL_cert_file => $plugin->opts->certificate,
	SSL_key_file => $plugin->opts->key
});

my $scriptContent = read_file($plugin->opts->script);

my $input = {
	"path" => $plugin->opts->executable,
	"args" => $plugin->opts->executableArg,
	"stdin" => $scriptContent,
	"timeout" => $plugin->opts->timeout,
	"scriptarguments" => \@ARGV,
};

if( -e $plugin->opts->script.".minisig") {
	my $signatureContent = read_file($plugin->opts->script.".minisig");
	$input->{"stdinsignature"} = $signatureContent;
}

alarm $plugin->opts->timeout;

my $response = $http->request('POST', "https://".urlize($username).":".urlize($password).'@'.$plugin->opts->hostname . ":" . $plugin->opts->port."/v1/runscriptstdin",{
	"headers" => {
		'Content-Type'=> 'application/json',
	},
	"content" => encode_json $input
});

if ($response->{status} ne 200) {
	print $response->{message}."\n".$response->{content}."\n";
	exit UNKNOWN
}

my $response_object = decode_json $response->{content};

print $response_object->{'output'};

my $exitcode = $response_object->{'exitcode'};
$exitcode = UNKNOWN if $exitcode > UNKNOWN;

exit $exitcode;

sub read_file {
	(my $filename) = @_;
	my $output = "";

	open(FH, '<', $filename) or die $!;

	while(<FH>){
		$output .= $_;
	}

	close(FH);
	return $output;
}

sub urlize {
	my ($rv) = @_;
	$rv =~ s/([^A-Za-z0-9])/sprintf("%%%2.2X", ord($1))/ge;
	return $rv;
}
