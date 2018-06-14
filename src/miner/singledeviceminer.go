package miner

import (
	"time"

	"github.com/robvanmieghem/go-opencl/cl"

	"github.com/modeneis/uminer/src/model"
	"github.com/modeneis/uminer/src/util"
)

func (miner *singleDeviceMiner) mine() {
	log.WithField("MinerID", miner.MinerID).Info("Initializing", miner.ClDevice.Type(), "-", miner.ClDevice.Name())

	context, err := cl.CreateContext([]*cl.Device{miner.ClDevice})
	if err != nil {
		log.WithError(err).WithField("MinerID", miner.MinerID).Fatalln("Failed to CreateContext")
	}
	defer context.Release()

	commandQueue, err := context.CreateCommandQueue(miner.ClDevice, 0)
	if err != nil {
		log.WithError(err).WithField("MinerID", miner.MinerID).Fatalln("Failed to CreateCommandQueue")
	}
	defer commandQueue.Release()

	program, err := context.CreateProgramWithSource([]string{kernelSource})
	if err != nil {
		log.WithError(err).WithField("MinerID", miner.MinerID).Fatalln("Failed to CreateProgramWithSource: ")
	}
	defer program.Release()

	err = program.BuildProgram([]*cl.Device{miner.ClDevice}, "")
	if err != nil {
		log.WithError(err).WithField("MinerID", miner.MinerID).Fatalln("Failed to BuildProgram: ")
	}

	kernel, err := program.CreateKernel("nonceGrind")
	if err != nil {
		log.WithError(err).WithField("MinerID", miner.MinerID).Fatalln("Failed to CreateKernel: ")
	}
	defer kernel.Release()

	blockHeaderObj := util.CreateEmptyBuffer(context, cl.MemReadOnly, 80)
	defer blockHeaderObj.Release()
	err = kernel.SetArgBuffer(0, blockHeaderObj)
	if err != nil {
		log.WithError(err).WithField("blockHeaderObj", blockHeaderObj).WithField("MinerID", miner.MinerID).Fatalln("Failed to blockHeaderObj kernel SetArgBuffer")
	}

	nonceOutObj := util.CreateEmptyBuffer(context, cl.MemReadWrite, 8)
	defer nonceOutObj.Release()
	err = kernel.SetArgBuffer(1, nonceOutObj)
	if err != nil {
		log.WithError(err).WithField("nonceOutObj", nonceOutObj).WithField("MinerID", miner.MinerID).Fatalln("Failed to nonceOutObj kernel SetArgBuffer", miner.MinerID, "-")
	}

	localItemSize, err := kernel.WorkGroupSize(miner.ClDevice)
	if err != nil {
		log.WithError(err).WithField("MinerID", miner.MinerID).Fatalln("ERROR: ", miner.MinerID, "- WorkGroupSize failed -")
	}

	log.WithField("MinerID", miner.MinerID).Debug("Global item size:", miner.GlobalItemSize, "(Intensity", miner.Intensity, ")", "- Local item size:", localItemSize)

	log.WithField("MinerID", miner.MinerID).Debug("Initialized ", miner.ClDevice.Type(), "-", miner.ClDevice.Name())

	//nonceOut := make([]byte, 8, 8)
	nonceOut := make([]byte, 8)
	if _, err = commandQueue.EnqueueWriteBufferByte(nonceOutObj, true, 0, nonceOut, nil); err != nil {
		log.WithError(err).WithField("MinerID", miner.MinerID).Fatalln("Failed to commandQueue.EnqueueWriteBufferByte")
	}
	for {
		start := time.Now()
		var work *miningWork
		var continueMining bool
		select {
		case work, continueMining = <-miner.miningWorkChannel:
		default:
			log.WithField("MinerID", miner.MinerID).Debug("No work ready")
			work, continueMining = <-miner.miningWorkChannel
			log.WithField("MinerID", miner.MinerID).Debug("Continuing")
		}
		if !continueMining {
			log.WithField("MinerID", miner.MinerID).Debug("Halting miner ")
			break
		}
		//Copy input to kernel args
		if _, err = commandQueue.EnqueueWriteBufferByte(blockHeaderObj, true, 0, work.Header, nil); err != nil {
			log.WithError(err).WithField("MinerID", miner.MinerID).Fatalln("Failed to Copy input to kernel args")
		}

		//Run the kernel
		if _, err = commandQueue.EnqueueNDRangeKernel(kernel, []int{work.Offset}, []int{miner.GlobalItemSize}, []int{localItemSize}, nil); err != nil {
			log.WithError(err).WithField("MinerID", miner.MinerID).Fatalln("Failed to Run the kernel")
		}
		//Get output
		if _, err = commandQueue.EnqueueReadBufferByte(nonceOutObj, true, 0, nonceOut, nil); err != nil {
			log.WithError(err).WithField("MinerID", miner.MinerID).Fatalln("Failed to Get output")
		}
		//Check if match found
		if nonceOut[0] != 0 || nonceOut[1] != 0 || nonceOut[2] != 0 || nonceOut[3] != 0 || nonceOut[4] != 0 || nonceOut[5] != 0 || nonceOut[6] != 0 || nonceOut[7] != 0 {
			log.WithField("MinerID", miner.MinerID).Info("Yay, solution found!")

			// Copy nonce to a new header.
			header := append([]byte(nil), work.Header...)
			for i := 0; i < 8; i++ {
				header[i+32] = nonceOut[i]
			}
			go func() {
				if e := miner.Client.SubmitHeader(header, work.Job); e != nil {
					log.WithError(e).WithField("MinerID", miner.MinerID).Println("Error submitting solution")
				}
			}()

			//Clear the output since it is dirty now
			//nonceOut = make([]byte, 8, 8)
			nonceOut = make([]byte, 8)
			if _, err = commandQueue.EnqueueWriteBufferByte(nonceOutObj, true, 0, nonceOut, nil); err != nil {
				log.WithError(err).WithField("MinerID", miner.MinerID).Fatalln("Failed to Clear the output ")
			}
		}

		hashRate := float64(miner.GlobalItemSize) / (time.Since(start).Seconds() * 1000000)
		miner.HashRateReports <- &model.HashRateReport{MinerID: miner.MinerID, HashRate: hashRate}
	}

}
