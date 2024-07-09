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


mkdir -p $SCRATCH/firecrest/$SLURM_JOBID

# Get node IP
node_name=$(scontrol show hostname $SLURM_JOB_NODELIST)
node_ip=$(getent hosts $node_name | awk '{ print $1 }')

# Log the node IP
echo "Node name: $node_name"
echo "Node IP: $node_ip"
echo $node_ip > $SCRATCH/firecrest/$SLURM_JOBID/node_ip.txt
echo "IP address written to $SCRATCH/firecrest/$SLURM_JOBID/node_ip.txt"
