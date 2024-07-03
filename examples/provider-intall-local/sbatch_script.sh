#!/bin/bash -l
#SBATCH --job-name="coder-job"
#SBATCH --mail-type=ALL
#SBATCH --mail-user="nicolotafta@gmail.com"
#SBATCH --time=00:05:00
#SBATCH --nodes=1
#SBATCH --ntasks-per-core=1
#SBATCH --ntasks-per-node=1
#SBATCH --cpus-per-task=12
#SBATCH --partition=debug
#SBATCH --constraint=gpu


