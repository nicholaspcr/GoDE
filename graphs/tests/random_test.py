import numpy as np

def correct_to_01(X, epsilon=1.0e-10):
    X[np.logical_and(X < 0, X >= 0 - epsilon)] = 0
    X[np.logical_and(X > 1, X <= 1 + epsilon)] = 1
    return X

def _shape_mixed(x, A=5.0, alpha=1.0):
    aux = 2.0 * A * np.pi
    ret = np.power(1.0 - x - (np.cos(aux * x + 0.5 * np.pi) / aux), alpha)
    return correct_to_01(ret)


def _calculate(x, s, h):
        return x[:, -1][:, None] + s * np.column_stack(h)


# variables

n_obj = 3
n_var = 8

y = np.array([0.1, 0.2, 0.3 ,0.4, -0.1, -0.2, -0.3, -0.4])
s = np.arange(2, 2 * n_obj + 1, 2)
h = np.array([1,2,3])


print(y,s,h)
print(y[:, -1][:,None])
# print(_calculate(y,s, h))

