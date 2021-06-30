package network::dns;

=head1 NAME

network::dns

=head1 DESCRIPTION

Do not manually modify this file.

This package was generated based on APIJSON/network/DNS.json,
validateJSON/network/DNS.json and validationFunctions/network/dns.pm.

=cut

use Data::Dumper;

use CyberAPIFunction;
use SFOS::Common::Utils;
use SFOS::Logging 'entity';

do '/_conf/csc/cscvalidation/Constants.pl';

our $entity = {
    "DNSFwdSelection" => {
        "datatype" => "STRING",
        "expectedvalues" => [
            "0",
            "1",
            "2",
            "3"
        ],
        "defaultvalue" => "0",
        "orm" => {"dbColumn" => "DNSFwdSelection"},
        "required" => "false"
    },
    "_entitydisplayname" => "DNS List",
    "dns3" => {
        "invalidinput" => [
            "MULTICAST",
            "BROADCAST",
            "UNSPECIFIED",
            "LINKLOCAL",
            "RESERVED"
        ],
        "datatype" => "IPADDRESS",
        "orm" => {"dbColumn" => "dns3"},
        "required" => "false"
    },
    "rdodhcpserver" => {
        "datatype" => "STRING",
        "expectedvalues" => [
            "0",
            "1",
            "2"
        ],
        "orm" => {"dbColumn" => "DNSServerFlag"},
        "required" => "false"
    },
    "dns2" => {
        "invalidinput" => [
            "MULTICAST",
            "BROADCAST",
            "UNSPECIFIED",
            "LINKLOCAL",
            "RESERVED"
        ],
        "datatype" => "IPADDRESS",
        "orm" => {"dbColumn" => "dns2"},
        "required" => "false"
    },
    "ipv6dns1" => {
        "invalidinput" => [
            "MULTICAST",
            "UNSPECIFIED",
            "LINKLOCAL",
            "IPV4MAPPED"
        ],
        "datatype" => "IPADDRESS6",
        "orm" => {"dbColumn" => "ipv6dns1"},
        "required" => "false"
    },
    "dns1" => {
        "invalidinput" => [
            "MULTICAST",
            "BROADCAST",
            "UNSPECIFIED",
            "LINKLOCAL",
            "RESERVED"
        ],
        "datatype" => "IPADDRESS",
        "orm" => {"dbColumn" => "dns1"},
        "required" => "false"
    },
    "ipv6rdodhcpserver" => {
        "datatype" => "STRING",
        "expectedvalues" => [
            "0",
            "1"
        ],
        "orm" => {"dbColumn" => "DNS6ServerFlag"},
        "required" => "false"
    },
    "displayname" => "DNS List",
    "ipv6dns3" => {
        "invalidinput" => [
            "MULTICAST",
            "UNSPECIFIED",
            "LINKLOCAL",
            "IPV4MAPPED"
        ],
        "datatype" => "IPADDRESS6",
        "orm" => {"dbColumn" => "ipv6dns3"},
        "required" => "false"
    },
    "ipv6dns2" => {
        "invalidinput" => [
            "MULTICAST",
            "UNSPECIFIED",
            "LINKLOCAL",
            "IPV4MAPPED"
        ],
        "datatype" => "IPADDRESS6",
        "orm" => {"dbColumn" => "ipv6dns2"},
        "required" => "false"
    },
    "validate" => "validateCheckForDuplicate"
};

our $EventProperties = {
    "READ" => {"ORM" => "false"},
    "DNS_CONFIGURATION" => {
        "readEntityName" => "DNS",
        "entityFilename" => "network::dns"
    }
};

our $LogProperties = {};

sub new {
    my $self = shift;
    my $data = shift;

    DEBUG('Input data for the object instantiation: ' . Dumper($data));
    DEBUG('Current state of entity: ' . Dumper($entity));

    if (!$data) {
        $data = {};
    }
    elsif (ref($data) eq 'HASH') {
        $data = getEntityJson($data);
    }
    else {
        ERROR("Parameter '\$data' must be a HASH reference.");
    }

    DEBUG('Entity is initialized to: ' . Dumper($entity));

    return bless $data, $self;
}

sub getReferenceBlock {
    my $self = shift;

    my $DbStructure = {
        "uniqueField" => "",
        "objectName" => "",
        "objectID" => "",
        "tableName" => ""
    };

    return $DbStructure;
}

sub getEntityJson {
    my $data = shift;

    if (exists &preProcess) {
        DEBUG('Calling preProcess function with data: ' . Dumper($data));
        preProcess($data);
    }

    # removing remaining left-over values from entity
    foreach my $key (keys %{$entity}) {
        if (ref($entity->{$key}) eq 'HASH' && defined $entity->{$key}{value}) {
            delete($entity->{$key}{value});
            DEBUG("The value of the following key is deleted from the entity: '$key'");
        }
    }

    # assigning values from the request to the entity
    foreach my $key (keys %{$data}) {
        $entity->{$key}{value} = $data->{$key};
        DEBUG("The value of the following key is overwritten in the entity: '$key'. New value: " . Dumper($data->{$key}));
    }

    my @trimKeys    = ('__newname');
    my $uniqueField = getReferenceBlock()->{uniqueField};
    if (defined $uniqueField && $uniqueField ne '') {
        push(@trimKeys, $uniqueField);
    }
    CyberAPIFunction::trimInputParameters($data, $entity, \@trimKeys);

    if (exists &logicalValidations) {
        DEBUG('Calling logicalValidation function with data: ' . Dumper($data));
        logicalValidations($data);
    }

    return $entity;
}

sub validateCheckForDuplicate{
	print "\n\n validateCheckForDuplicate Function Called*********";
	my $entityJson=shift;
	my @errorlist=();
	
	$dns1=$entityJson->{dns1}->{value};
	$dns2=$entityJson->{dns2}->{value};
	$dns3=$entityJson->{dns3}->{value};
	$ipv6dns1=$entityJson->{ipv6dns1}->{value};
	$ipv6dns2=$entityJson->{ipv6dns2}->{value};
	$ipv6dns3=$entityJson->{ipv6dns3}->{value};
	if($dns1 eq '' && $dns2 ne ''){
		$request->{dns1}=$dns2;
		$dns1=$dns2;
		$request->{dns2}="";
		$dns2="";		
	}
	if($dns2 eq '' && $dns3 ne ''){
		if($dns1 eq ''){
			$dns1=$dns3;
			$request->{dns1}=$dns3;
		}else{
			$dns2=$dns3;
			$request->{dns2}=$dns3;
		}		
		$request->{dns3}="";
		$dns3="";
	}
		
	if((defined $dns1 && defined $dns2) && (($dns1 eq $dns2 || $dns1 eq $dns3) || ($dns2 eq $dns3))){			
		push(@errorlist,"Duplicate DNS IPv4 Entries.");
		if($dns1 eq $dns2){
			$request->{dns1}=$dns2;
			$dns1=$dns2;
			$request->{dns2}=$dns3;
			$dns2=$dns3;
			$request->{dns3}="";
			$dns3="";
			if($dns1 eq $dns2){
				$request->{dns1}=$dns2;
				$dns1=$dns2;
				$request->{dns2}=$dns3;
				$dns2=$dns3;
				$request->{dns3}="";
				$dns3="";
			}
		}			
		if($dns1 eq $dns3){
			$request->{dns1}=$dns1;
			$dns1=$dns1;
			$request->{dns2}=$dns2;
			$dns2=$dns2;
			$request->{dns3}="";
			$dns3="";
		}
		if($dns2 eq $dns3){
			$request->{dns1}=$dns1;
			$dns1=$dns1;
			$request->{dns2}=$dns3;
			$dns2=$dns3;
			$request->{dns3}="";
			$dns3="";
		}
	}

	if($ipv6dns1 eq '' && $ipv6dns2 ne ''){
		$request->{ipv6dns1}=$ipv6dns2;
		$ipv6dns1=$ipv6dns2;
		$request->{ipv6dns2}="";
		$ipv6dns2="";		
	}
	if($ipv6dns2 eq '' && $ipv6dns3 ne ''){
		if($ipv6dns1 eq ''){
			$ipv6dns1=$ipv6dns3;
			$request->{ipv6dns1}=$ipv6dns3;
		}else{
			$ipv6dns2=$ipv6dns3;
			$request->{ipv6dns2}=$ipv6dns3;
		}		
		$request->{ipv6dns3}="";
		$ipv6dns3="";
	}
	
	if((defined $ipv6dns1 && defined $ipv6dns2) && (($ipv6dns1 eq $ipv6dns2 || $ipv6dns1 eq $ipv6dns3) || ($ipv6dns2 eq $ipv6dns3))){
		if($ipv6dns1 eq $ipv6dns2){
			$request->{ipv6dns1}=$ipv6dns2;
			$ipv6dns1=$ipv6dns2;
			$request->{ipv6dns2}=$ipv6dns3;
			$ipv6dns2=$ipv6dns3;
			$request->{ipv6dns3}="";
			$ipv6dns3="";
			if($ipv6dns1 eq $ipv6dns2){
				print "\n\n Sec Iff";
				$request->{ipv6dns1}=$ipv6dns2;
				$ipv6dns1=$ipv6dns2;
				$request->{ipv6dns2}=$ipv6dns3;
				$ipv6dns2=$ipv6dns3;
				$request->{ipv6dns3}="";
				$ipv6dns3="";
			}
		}			
		if($ipv6dns1 eq $ipv6dns3){
			$request->{ipv6dns1}=$ipv6dns1;
			$ipv6dns1=$ipv6dns1;
			$request->{ipv6dns2}=$ipv6dns2;
			$ipv6dns2=$ipv6dns2;
			$request->{ipv6dns3}="";
			$ipv6dns3="";
		}
		if($ipv6dns2 eq $ipv6dns3){
			$request->{ipv6dns1}=$ipv6dns1;
			$ipv6dns1=$ipv6dns1;
			$request->{ipv6dns2}=$ipv6dns3;
			$ipv6dns2=$ipv6dns3;
			$request->{ipv6dns3}="";
			$ipv6dns3="";
		}
		push(@errorlist,"Duplicate DNS IPv6 Entries.");
	}
	return \@errorlist;
}
sub logicalValidations{
	$request=shift;
	if($request->{dns1} eq ''){
		delete($entity->{dns1});
		delete($request->{dns1});
	}
	if($request->{dns2} eq ''){
		delete($entity->{dns2});
		delete($request->{dns2});
	}
	if($request->{dns3} eq ''){
		delete($entity->{dns3});
		delete($request->{dns3});
	}
	if($request->{ipv6dns1} eq ''){
		delete($entity->{ipv6dns1});
		delete($request->{ipv6dns1});
	}
	if($request->{ipv6dns2} eq ''){
		delete($entity->{ipv6dns2});
		delete($request->{ipv6dns2});
	}
	if($request->{ipv6dns3} eq ''){
		delete($entity->{ipv6dns3});
		delete($request->{ipv6dns3});
	}
}
1;