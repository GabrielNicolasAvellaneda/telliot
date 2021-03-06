// Copyright (c) The Tellor Authors.
// Licensed under the MIT License.

package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"github.com/tellor-io/telliot/pkg/db"
	"github.com/tellor-io/telliot/pkg/ops"
)

var GitTag string
var GitHash string

const versionMessage = `
    The official Tellor cli tool %s (%s)
    -----------------------------------------
	Website: https://tellor.io
	Github:  https://github.com/tellor-io/telliot
`

type VersionCmd struct {
}

func (cmd *VersionCmd) Run() error {
	fmt.Printf(versionMessage, GitTag, GitHash)
	return nil
}

type configPath string
type tokenCmd struct {
	Config  configPath `required:"" type:"existingfile" help:"path to config file"`
	Address string     `arg:""`
	Amount  string     `arg:""`
}

type transferCmd tokenCmd

func (c *transferCmd) Run() error {
	cfg, err := parseConfig(string(c.Config))
	if err != nil {
		return errors.Wrapf(err, "creating config")
	}

	logger, err := createLogger(cfg.Logger, cfg.LogLevel)
	if err != nil {
		return errors.Wrapf(err, "creating logger")
	}

	ctx := context.Background()
	client, contract, account, err := createTellorVariables(ctx, cfg)
	if err != nil {
		return errors.Wrapf(err, "creating tellor variables")
	}

	address := ETHAddress{}
	err = address.Set(c.Address)
	if err != nil {
		return errors.Wrapf(err, "parsing address argument")
	}
	amount := TRBAmount{}
	err = amount.Set(c.Amount)
	if err != nil {
		return errors.Wrapf(err, "parsing amount argument")
	}
	return ops.Transfer(ctx, logger, client, contract, account, address.addr, amount.Int)
}

type approveCmd tokenCmd

func (c *approveCmd) Run() error {
	cfg, err := parseConfig(string(c.Config))
	if err != nil {
		return errors.Wrapf(err, "creating config")
	}

	logger, err := createLogger(cfg.Logger, cfg.LogLevel)
	if err != nil {
		return errors.Wrapf(err, "creating logger")
	}

	ctx := context.Background()
	client, contract, account, err := createTellorVariables(ctx, cfg)
	if err != nil {
		return errors.Wrapf(err, "creating tellor variables")
	}

	address := ETHAddress{}
	err = address.Set(c.Address)
	if err != nil {
		return errors.Wrapf(err, "parsing address argument")
	}
	amount := TRBAmount{}
	err = amount.Set(c.Amount)
	if err != nil {
		return errors.Wrapf(err, "parsing amount argument")
	}
	return ops.Approve(ctx, logger, client, contract, account, address.addr, amount.Int)
}

type balanceCmd struct {
	Config  configPath `required:"" type:"existingfile" help:"path to config file"`
	Address string     `arg:"" optional:""`
}

func (b *balanceCmd) Run() error {
	cfg, err := parseConfig(string(b.Config))
	if err != nil {
		return errors.Wrapf(err, "creating config")
	}

	_, err = createLogger(cfg.Logger, cfg.LogLevel)
	if err != nil {
		return errors.Wrapf(err, "creating logger")
	}

	ctx := context.Background()
	client, contract, _, err := createTellorVariables(ctx, cfg)
	if err != nil {
		return errors.Wrapf(err, "creating tellor variables")
	}

	addr := ETHAddress{}
	if b.Address == "" {
		err = addr.Set(contract.Address.String())
		if err != nil {
			return errors.Wrapf(err, "parsing argument")
		}
	} else {
		err = addr.Set(b.Address)
		if err != nil {
			return errors.Wrapf(err, "parsing argument")
		}
	}
	return ops.Balance(ctx, client, contract.Getter, addr.addr)
}

type depositCmd struct {
	Config configPath `required:"" type:"existingfile" help:"path to config file"`
}

func (d depositCmd) Run() error {
	cfg, err := parseConfig(string(d.Config))
	if err != nil {
		return errors.Wrapf(err, "creating config")
	}

	logger, err := createLogger(cfg.Logger, cfg.LogLevel)
	if err != nil {
		return errors.Wrapf(err, "creating logger")
	}

	ctx := context.Background()
	client, contract, account, err := createTellorVariables(ctx, cfg)
	if err != nil {
		return errors.Wrapf(err, "creating tellor variables")
	}
	return ops.Deposit(ctx, logger, client, contract, account)
}

type withdrawCmd struct {
	Config configPath `required:"" type:"existingfile" help:"path to config file"`
}

func (w withdrawCmd) Run() error {
	cfg, err := parseConfig(string(w.Config))
	if err != nil {
		return errors.Wrapf(err, "creating config")
	}

	logger, err := createLogger(cfg.Logger, cfg.LogLevel)
	if err != nil {
		return errors.Wrapf(err, "creating logger")
	}

	ctx := context.Background()
	client, contract, account, err := createTellorVariables(ctx, cfg)
	if err != nil {
		return errors.Wrapf(err, "creating tellor variables")
	}
	return ops.WithdrawStake(ctx, logger, client, contract, account)
}

type requestCmd struct {
	Config configPath `required:"" type:"existingfile" help:"path to config file"`
}

func (r requestCmd) Run() error {
	cfg, err := parseConfig(string(r.Config))
	if err != nil {
		return errors.Wrapf(err, "creating config")
	}

	logger, err := createLogger(cfg.Logger, cfg.LogLevel)
	if err != nil {
		return errors.Wrapf(err, "creating logger")
	}

	ctx := context.Background()
	client, contract, account, err := createTellorVariables(ctx, cfg)
	if err != nil {
		return errors.Wrapf(err, "creating tellor variables")
	}
	return ops.RequestStakingWithdraw(ctx, logger, client, contract, account)
}

type statusCmd struct {
	Config configPath `required:"" type:"existingfile" help:"path to config file"`
}

func (s statusCmd) Run() error {
	cfg, err := parseConfig(string(s.Config))
	if err != nil {
		return errors.Wrapf(err, "creating config")
	}

	logger, err := createLogger(cfg.Logger, cfg.LogLevel)
	if err != nil {
		return errors.Wrapf(err, "creating logger")
	}

	ctx := context.Background()
	client, contract, account, err := createTellorVariables(ctx, cfg)
	if err != nil {
		return errors.Wrapf(err, "creating tellor variables")
	}
	return ops.ShowStatus(ctx, logger, client, contract, account)
}

type newDisputeCmd struct {
	Config     configPath `required:"" type:"existingfile" help:"path to config file"`
	requestId  string     `arg:"" required:"" help:"the request id to dispute it"`
	timestamp  string     `arg:"" required:"" help:"the submitted timestamp to dispute"`
	minerIndex string     `arg:"" required:"" help:"the miner index to dispute"`
}

func (n newDisputeCmd) Run() error {
	cfg, err := parseConfig(string(n.Config))
	if err != nil {
		return errors.Wrapf(err, "creating config")
	}

	_, err = createLogger(cfg.Logger, cfg.LogLevel)
	if err != nil {
		return errors.Wrapf(err, "creating logger")
	}

	ctx := context.Background()
	client, contract, account, err := createTellorVariables(ctx, cfg)
	if err != nil {
		return errors.Wrapf(err, "creating tellor variables")
	}

	requestID := EthereumInt{}
	err = requestID.Set(n.requestId)
	if err != nil {
		return errors.Wrapf(err, "parsing argument")
	}
	timestamp := EthereumInt{}
	err = timestamp.Set(n.timestamp)
	if err != nil {
		return errors.Wrapf(err, "parsing argument")
	}
	minerIndex := EthereumInt{}
	err = minerIndex.Set(n.minerIndex)
	if err != nil {
		return errors.Wrapf(err, "parsing argument")
	}
	return ops.Dispute(ctx, client, contract, account, requestID.Int, timestamp.Int, minerIndex.Int)
}

type voteCmd struct {
	Config    configPath `required:"" type:"existingfile" help:"path to config file"`
	disputeId string     `arg:"" required:"" help:"the dispute id"`
	support   bool       `arg:"" required:"" help:"true or false"`
}

func (v voteCmd) Run() error {
	cfg, err := parseConfig(string(v.Config))
	if err != nil {
		return errors.Wrapf(err, "creating config")
	}

	_, err = createLogger(cfg.Logger, cfg.LogLevel)
	if err != nil {
		return errors.Wrapf(err, "creating logger")
	}

	ctx := context.Background()
	client, contract, account, err := createTellorVariables(ctx, cfg)
	if err != nil {
		return errors.Wrapf(err, "creating tellor variables")
	}

	disputeID := EthereumInt{}
	err = disputeID.Set(v.disputeId)
	if err != nil {
		return errors.Wrapf(err, "parsing argument")
	}
	return ops.Vote(ctx, client, contract, account, disputeID.Int, v.support)
}

type showCmd struct {
	Config configPath `required:"" type:"existingfile" help:"path to config file"`
}

func (s showCmd) Run() error {
	cfg, err := parseConfig(string(s.Config))
	if err != nil {
		return errors.Wrapf(err, "creating config")
	}

	logger, err := createLogger(cfg.Logger, cfg.LogLevel)
	if err != nil {
		return errors.Wrapf(err, "creating logger")
	}

	ctx := context.Background()
	client, contract, account, err := createTellorVariables(ctx, cfg)
	if err != nil {
		return errors.Wrapf(err, "creating tellor variables")
	}
	return ops.List(ctx, logger, client, contract, account)
}

type dataserverCmd struct {
	Config configPath `required:"" type:"existingfile" help:"path to config file"`
}

func (d dataserverCmd) Run() error {
	cfg, err := parseConfig(string(d.Config))
	if err != nil {
		return errors.Wrapf(err, "creating config")
	}

	logger, err := createLogger(cfg.Logger, cfg.LogLevel)
	if err != nil {
		return errors.Wrapf(err, "creating logger")
	}

	ctx := context.Background()
	client, contract, account, err := createTellorVariables(ctx, cfg)
	if err != nil {
		return errors.Wrapf(err, "creating tellor variables")
	}

	// Create os kill sig listener.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	var ds *ops.DataServerOps
	DB, err := migrateAndOpenDB(cfg)
	if err != nil {
		return errors.Wrapf(err, "initializing database")
	}
	proxy, err := createProxy(cfg, DB)
	if err != nil {
		return errors.Wrapf(err, "initializing proxy")
	}
	ch := make(chan os.Signal)
	ds, err = ops.CreateDataServerOps(ctx, logger, cfg, DB, &proxy, client, contract, account, ch)
	if err != nil {
		return errors.Wrapf(err, "creating data server")
	}
	// Start and wait for it to be ready
	if err := ds.Start(ctx); err != nil {
		return errors.Wrapf(err, "starting data server")
	}
	<-ds.Ready()

	// Wait for kill sig.
	<-c
	// Notify exit channels.
	ch <- os.Interrupt

	cnt := 0
	start := time.Now()
	for {
		cnt++
		dsStopped := false

		if ds != nil {
			dsStopped = !ds.Running
		} else {
			dsStopped = true
		}

		if !dsStopped && cnt > 60 {
			level.Warn(logger).Log("msg", "taking longer than expected to stop operations", "waited", time.Since(start))
		} else if dsStopped {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	level.Info(logger).Log("msg", "main shutdown complete")
	return nil
}

type mineCmd struct {
	Config configPath `required:"" type:"existingfile" help:"path to config file"`
}

func (m mineCmd) Run() error {
	cfg, err := parseConfig(string(m.Config))
	if err != nil {
		return errors.Wrapf(err, "creating config")
	}

	logger, err := createLogger(cfg.Logger, cfg.LogLevel)
	if err != nil {
		return errors.Wrapf(err, "creating logger")
	}

	ctx := context.Background()
	client, contract, account, err := createTellorVariables(ctx, cfg)
	if err != nil {
		return errors.Wrapf(err, "creating tellor variables")
	}
	// Create os kill sig listener.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	exitChannels := make([]*chan os.Signal, 0)

	var ds *ops.DataServerOps
	DB, err := migrateAndOpenDB(cfg)
	if err != nil {
		return errors.Wrapf(err, "initializing database")
	}
	proxy, err := createProxy(cfg, DB)
	if err != nil {
		return errors.Wrapf(err, "initializing proxy")
	}
	if !cfg.EnablePoolWorker {
		if cfg.Mine.RemoteDBHost == "" {
			ch := make(chan os.Signal)
			exitChannels = append(exitChannels, &ch)

			ds, err = ops.CreateDataServerOps(ctx, logger, cfg, DB, &proxy, client, contract, account, ch)
			if err != nil {
				return errors.Wrapf(err, "creating data server")
			}
			// Start and wait for it to be ready.
			if err := ds.Start(ctx); err != nil {
				return errors.Wrapf(err, "starting data server")
			}
			<-ds.Ready()
		}
	}
	// Start miner
	v, err := DB.Get(db.DisputeStatusKey)
	if err != nil {
		level.Warn(logger).Log("msg", "getting dispute status. Check if staked")
	}
	status, _ := hexutil.DecodeBig(string(v))
	if status.Cmp(big.NewInt(1)) != 0 {
		return errors.New("miner is not able to mine with current status")
	}
	ch := make(chan os.Signal)
	exitChannels = append(exitChannels, &ch)
	miner, err := ops.CreateMiningManager(logger, ch, cfg, proxy, contract, account)
	if err != nil {
		return errors.Wrapf(err, "creating miner")
	}
	go func() {
		miner.Start(ctx)
	}()

	// Wait for kill sig.
	<-c
	// Then notify exit channels.
	for _, ch := range exitChannels {
		*ch <- os.Interrupt
	}
	cnt := 0
	start := time.Now()
	for {
		cnt++
		dsStopped := false
		minerStopped := false

		if ds != nil {
			dsStopped = !ds.Running
		} else {
			dsStopped = true
		}

		if miner != nil {
			minerStopped = !miner.Running
		} else {
			minerStopped = true
		}

		if !dsStopped && !minerStopped && cnt > 60 {
			level.Warn(logger).Log("msg", "taking longer than expected to stop operations", "waited", time.Since(start))
		} else if dsStopped && minerStopped {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	level.Info(logger).Log("msg", "main shutdown complete")
	return nil
}
