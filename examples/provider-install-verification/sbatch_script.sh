#!/bin/bash -l
#SBATCH --job-name="job-test"
#SBATCH --mail-type=ALL
#SBATCH --mail-user="nicolotafta@gmail.com"
#SBATCH --time=00:01:00
#SBATCH --nodes=1
#SBATCH --ntasks-per-core=2
#SBATCH --ntasks-per-node=6
#SBATCH --cpus-per-task=1
#SBATCH --partition=debug
#SBATCH --constraint=gpu


