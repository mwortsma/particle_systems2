import matplotlib.pyplot as plt
from matplotlib import pylab
import numpy as np

def get_keys(distributions):
	for d1 in distributions:
		for d2 in distributions:
			for k in d1:
				if k not in d2: d2[k] = 0
	return sorted(distributions[0].keys())

def plot_path(distributions, labels, show, save):
	keys = get_keys(distributions)
	for i in range(len(distributions)):
		d = distributions[i]
		plt.plot([d[k] for k in keys], label=labels[i])
		print [d[k] for k in keys]
	for i in range(len(distributions)):
		for j in range(i+1,len(distributions)):
			print "Error between " + str(labels[i]) + " " + str(labels[j]) + ":"
			d1 = distributions[i]
			d2 = distributions[j]
			print sum([abs(d1[k]-d2[k]) for k in keys])
			print max([abs(d1[k]-d2[k]) for k in keys])
	plt.legend(loc=2)
	plt.ylabel("Probability")
	if show:
		plt.show()
	if save and save != "":
		plt.savefig(save)


def plot_time(distributions, labels, show, save):
	for i in range(len(distributions)):
		d = distributions[i]
		k = d['K']
		arr = np.array(d['Distr'])
		for j in range(k-1):
			plt.plot(np.arange(0, len(arr[:,j])*d['Dt'],d['Dt']), arr[:,j],
				label=(labels[i]+" P(X="+str(j))+")")
	plt.legend(loc=2)
	plt.xlabel("Time")
	if show:
		plt.show()
	if save and save != "":
		plt.savefig(save)


def plot_time_full(distributions, labels, show, save):
	for i in range(len(distributions)):
		d = distributions[i]
		k = d['K']
		arr = np.array(d['Distr'])
		for j in range(k):
			plt.plot(np.arange(0, len(arr[:,j])*d['Dt'],d['Dt']), arr[:,j],
				label=(labels[i]+" P(X="+str(j))+")")

	plt.legend(loc=2)
	if show:
		plt.show()
	if save and save != "":
		plt.savefig(save)
