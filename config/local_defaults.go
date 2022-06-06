// Copyright (C) 2019-2022 Algorand, Inc.
// This file is part of go-algorand
//
// go-algorand is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// go-algorand is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with go-algorand.  If not, see <https://www.gnu.org/licenses/>.

// This file was auto generated by ./config/defaultsGenerator/defaultsGenerator.go, and SHOULD NOT BE MODIFIED in any way
// If you want to make changes to this file, make the corresponding changes to Local in localTemplate.go and run "go generate".

package config

var defaultLocal = Local{
	Version:                                    22,
	AccountUpdatesStatsInterval:                5000000000,
	AccountsRebuildSynchronousMode:             1,
	AgreementIncomingBundlesQueueLength:        7,
	AgreementIncomingProposalsQueueLength:      25,
	AgreementIncomingVotesQueueLength:          10000,
	AnnounceParticipationKey:                   true,
	Archival:                                   false,
	BaseLoggerDebugLevel:                       4,
	BlockServiceCustomFallbackEndpoints:        "",
	BroadcastConnectionsLimit:                  -1,
	CadaverSizeTarget:                          1073741824,
	CatchpointFileHistoryLength:                365,
	CatchpointInterval:                         10000,
	CatchpointTracking:                         0,
	CatchupBlockDownloadRetryAttempts:          1000,
	CatchupBlockValidateMode:                   0,
	CatchupFailurePeerRefreshRate:              10,
	CatchupGossipBlockFetchTimeoutSec:          4,
	CatchupHTTPBlockFetchTimeoutSec:            4,
	CatchupLedgerDownloadRetryAttempts:         50,
	CatchupParallelBlocks:                      16,
	ConnectionsRateLimitingCount:               60,
	ConnectionsRateLimitingWindowSeconds:       1,
	DNSBootstrapID:                             "<network>.algorand.network",
	DNSSecurityFlags:                           1,
	DeadlockDetection:                          0,
	DeadlockDetectionThreshold:                 30,
	DisableLocalhostConnectionRateLimit:        true,
	DisableNetworking:                          false,
	DisableOutgoingConnectionThrottling:        false,
	EnableAccountUpdatesStats:                  false,
	EnableAgreementReporting:                   false,
	EnableAgreementTimeMetrics:                 false,
	EnableAssembleStats:                        false,
	EnableBlockService:                         false,
	EnableBlockServiceFallbackToArchiver:       true,
	EnableCatchupFromArchiveServers:            false,
	EnableDeveloperAPI:                         false,
	EnableGossipBlockService:                   true,
	EnableIncomingMessageFilter:                false,
	EnableLedgerService:                        false,
	EnableMetricReporting:                      false,
	EnableOutgoingNetworkMessageFiltering:      true,
	EnablePingHandler:                          true,
	EnableProcessBlockStats:                    false,
	EnableProfiler:                             false,
	EnableRequestLogger:                        false,
	EnableRuntimeMetrics:                       false,
	EnableTopAccountsReporting:                 false,
	EnableVerbosedTransactionSyncLogging:       false,
	EndpointAddress:                            "127.0.0.1:0",
	FallbackDNSResolverAddress:                 "",
	ForceFetchTransactions:                     false,
	ForceRelayMessages:                         false,
	GossipFanout:                               4,
	IncomingConnectionsLimit:                   800,
	IncomingMessageFilterBucketCount:           5,
	IncomingMessageFilterBucketSize:            512,
	IsIndexerActive:                            false,
	LedgerSynchronousMode:                      2,
	LogArchiveMaxAge:                           "",
	LogArchiveName:                             "node.archive.log",
	LogSizeLimit:                               1073741824,
	MaxAPIResourcesPerAccount:                  100000,
	MaxCatchpointDownloadDuration:              7200000000000,
	MaxConnectionsPerIP:                        30,
	MinCatchpointFileDownloadBytesPerSecond:    20480,
	NetAddress:                                 "",
	NetworkMessageTraceServer:                  "",
	NetworkProtocolVersion:                     "",
	NodeExporterListenAddress:                  ":9100",
	NodeExporterPath:                           "./node_exporter",
	OptimizeAccountsDatabaseOnStartup:          false,
	OutgoingMessageFilterBucketCount:           3,
	OutgoingMessageFilterBucketSize:            128,
	ParticipationKeysRefreshInterval:           60000000000,
	PeerConnectionsUpdateInterval:              3600,
	PeerPingPeriodSeconds:                      0,
	PriorityPeers:                              map[string]bool{},
	ProposalAssemblyTime:                       250000000,
	PublicAddress:                              "",
	ReconnectTime:                              60000000000,
	ReservedFDs:                                256,
	RestConnectionsHardLimit:                   2048,
	RestConnectionsSoftLimit:                   1024,
	RestReadTimeoutSeconds:                     15,
	RestWriteTimeoutSeconds:                    120,
	RunHosted:                                  false,
	SuggestedFeeBlockHistory:                   3,
	SuggestedFeeSlidingWindowSize:              50,
	TLSCertFile:                                "",
	TLSKeyFile:                                 "",
	TelemetryToLog:                             true,
	TransactionSyncDataExchangeRate:            0,
	TransactionSyncSignificantMessageThreshold: 0,
	TxPoolExponentialIncreaseFactor:            2,
	TxPoolSize:                                 15000,
	TxSyncIntervalSeconds:                      60,
	TxSyncServeResponseSize:                    1000000,
	TxSyncTimeoutSeconds:                       30,
	UseXForwardedForAddressField:               "",
	VerifiedTranscationsCacheSize:              30000,
}
