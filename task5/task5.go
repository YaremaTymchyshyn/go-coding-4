package main

import (
	"errors"
	"fmt"
	"math"
)

type Matrix struct {
	rows, cols int
	data       [][]float64
}

type SortableMatrix interface {
	Sort()
}

type RowSortedMatrix struct {
	Matrix
}

func (r RowSortedMatrix) Sort() {
	for i := 0; i < len(r.data)-1; i++ {
		for j := 0; j < len(r.data)-i-1; j++ {
			if !lessThan(r.data[j], r.data[j+1]) {
				r.data[j], r.data[j+1] = r.data[j+1], r.data[j]
			}
		}
	}
}

func lessThan(row1, row2 []float64) bool {
	for i := 0; i < len(row1); i++ {
		if row1[i] != row2[i] {
			return row1[i] < row2[i]
		}
	}
	return false
}

type ColSortedMatrix struct {
	Matrix
}

func (c ColSortedMatrix) Sort() {
	for i := 0; i < c.cols-1; i++ {
		for j := 0; j < c.cols-i-1; j++ {
			if !lessThanCol(c.data, j, j+1) {
				swapCols(c.data, j, j+1)
			}
		}
	}
}

func swapCols(matrix [][]float64, col1, col2 int) {
	for i := 0; i < len(matrix); i++ {
		matrix[i][col1], matrix[i][col2] = matrix[i][col2], matrix[i][col1]
	}
}

func lessThanCol(matrix [][]float64, col1, col2 int) bool {
	for i := 0; i < len(matrix); i++ {
		if matrix[i][col1] != matrix[i][col2] {
			return matrix[i][col1] < matrix[i][col2]
		}
	}
	return false
}

func NewMatrix(rows, cols int) Matrix {
	data := make([][]float64, rows)
	for i := range data {
		data[i] = make([]float64, cols)
	}
	return Matrix{rows, cols, data}
}

func (m Matrix) String() string {
	res := ""
	for _, row := range m.data {
		for _, val := range row {
			res += fmt.Sprintf("%8.2f ", val)
		}
		res += "\n"
	}
	return res
}

func InputMatrix(rows, cols int) Matrix {
	matrix := NewMatrix(rows, cols)
	fmt.Println("Enter matrix elements row by row:")
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			fmt.Scan(&matrix.data[i][j])
		}
	}
	return matrix
}

func (m Matrix) Add(n Matrix) (Matrix, error) {
	if m.rows != n.rows || m.cols != n.cols {
		return Matrix{}, errors.New("matrices dimensions do not match")
	}
	result := NewMatrix(m.rows, m.cols)
	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			result.data[i][j] = m.data[i][j] + n.data[i][j]
		}
	}
	return result, nil
}

func (m Matrix) Subtract(n Matrix) (Matrix, error) {
	if m.rows != n.rows || m.cols != n.cols {
		return Matrix{}, errors.New("matrices dimensions do not match")
	}
	result := NewMatrix(m.rows, m.cols)
	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			result.data[i][j] = m.data[i][j] - n.data[i][j]
		}
	}
	return result, nil
}

func (m Matrix) Multiply(n Matrix) (Matrix, error) {
	if m.cols != n.rows {
		return Matrix{}, errors.New("matrices cannot be multiplied")
	}
	result := NewMatrix(m.rows, n.cols)
	for i := 0; i < m.rows; i++ {
		for j := 0; j < n.cols; j++ {
			for k := 0; k < m.cols; k++ {
				result.data[i][j] += m.data[i][k] * n.data[k][j]
			}
		}
	}
	return result, nil
}

func (m Matrix) Transpose() Matrix {
	result := NewMatrix(m.cols, m.rows)
	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			result.data[j][i] = m.data[i][j]
		}
	}
	return result
}

func (m Matrix) Determinant() (float64, error) {
	if m.rows != m.cols {
		return 0, errors.New("determinant is only defined for square matrices")
	}
	return determinant(m.data), nil
}

func determinant(mat [][]float64) float64 {
	n := len(mat)
	if n == 1 {
		return mat[0][0]
	} else if n == 2 {
		return mat[0][0]*mat[1][1] - mat[0][1]*mat[1][0]
	}

	var det float64
	for p := 0; p < n; p++ {
		subMatrix := make([][]float64, n-1)
		for i := range subMatrix {
			subMatrix[i] = make([]float64, n-1)
		}
		for i := 1; i < n; i++ {
			for j, col := 0, 0; j < n; j++ {
				if j != p {
					subMatrix[i-1][col] = mat[i][j]
					col++
				}
			}
		}
		det += math.Pow(-1, float64(p)) * mat[0][p] * determinant(subMatrix)
	}
	return det
}

func (m Matrix) Inverse() (Matrix, error) {
	if m.rows != m.cols {
		return Matrix{}, errors.New("only square matrices can be inverted")
	}
	det, err := m.Determinant()
	if err != nil || det == 0 {
		return Matrix{}, errors.New("matrix cannot be inverted")
	}
	cofactorMatrix := m.Cofactor()
	adjugate := cofactorMatrix.Transpose()
	inverseMatrix := NewMatrix(m.rows, m.cols)
	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			inverseMatrix.data[i][j] = adjugate.data[i][j] / det
		}
	}
	return inverseMatrix, nil
}

func (m Matrix) Cofactor() Matrix {
	cofactorMatrix := NewMatrix(m.rows, m.cols)
	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			minor := m.Minor(i, j)
			cofactorMatrix.data[i][j] = math.Pow(-1, float64(i+j)) * determinant(minor)
		}
	}
	return cofactorMatrix
}

func (m Matrix) Minor(row, col int) [][]float64 {
	minor := make([][]float64, m.rows-1)
	for i := range minor {
		minor[i] = make([]float64, m.cols-1)
	}

	for i, minorRow := 0, 0; i < m.rows; i++ {
		if i == row {
			continue
		}
		for j, minorCol := 0, 0; j < m.cols; j++ {
			if j == col {
				continue
			}
			minor[minorRow][minorCol] = m.data[i][j]
			minorCol++
		}
		minorRow++
	}
	return minor
}

func (m Matrix) SolveLinearEquations() ([]float64, error) {
	if m.cols != m.rows+1 {
		return nil, errors.New("augmented matrix must have one more column than rows")
	}

	for i := 0; i < m.rows; i++ {
		maxRow := i
		for k := i + 1; k < m.rows; k++ {
			if abs(m.data[k][i]) > abs(m.data[maxRow][i]) {
				maxRow = k
			}
		}
		m.data[i], m.data[maxRow] = m.data[maxRow], m.data[i]

		for k := i + 1; k < m.rows; k++ {
			factor := m.data[k][i] / m.data[i][i]
			for j := i; j < m.cols; j++ {
				m.data[k][j] -= factor * m.data[i][j]
			}
		}
	}

	x := make([]float64, m.rows)
	for i := m.rows - 1; i >= 0; i-- {
		x[i] = m.data[i][m.cols-1] / m.data[i][i]
		for k := i - 1; k >= 0; k-- {
			m.data[k][m.cols-1] -= m.data[k][i] * x[i]
		}
	}

	return x, nil
}

func abs(a float64) float64 {
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	fmt.Println("Matrix Operations Console")
	fmt.Println("Choose an operation:")
	fmt.Println("1 - Add")
	fmt.Println("2 - Subtract")
	fmt.Println("3 - Multiply")
	fmt.Println("4 - Transpose")
	fmt.Println("5 - Determinant")
	fmt.Println("6 - Inverse")
	fmt.Println("7 - Sort Rows")
	fmt.Println("8 - Sort Columns")
	fmt.Println("9 - Solve system of linear equations")

	var choice int
	fmt.Scan(&choice)

	switch choice {
	case 1, 2, 3:
		fmt.Println("Enter dimensions for matrix 1 (rows and columns):")
		var rows1, cols1 int
		fmt.Scan(&rows1, &cols1)
		matrix1 := InputMatrix(rows1, cols1)

		fmt.Println("Enter dimensions for matrix 2 (rows and columns):")
		var rows2, cols2 int
		fmt.Scan(&rows2, &cols2)
		matrix2 := InputMatrix(rows2, cols2)

		var result Matrix
		var err error
		if choice == 1 {
			result, err = matrix1.Add(matrix2)
		} else if choice == 2 {
			result, err = matrix1.Subtract(matrix2)
		} else {
			result, err = matrix1.Multiply(matrix2)
		}
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Result:")
			fmt.Println(result)
		}
	case 4:
		fmt.Println("Enter dimensions for the matrix (rows and columns):")
		var rows, cols int
		fmt.Scan(&rows, &cols)
		matrix := InputMatrix(rows, cols)
		result := matrix.Transpose()
		fmt.Println("Transposed matrix:")
		fmt.Println(result)
	case 5:
		fmt.Println("Enter dimensions for the matrix (rows and columns):")
		var rows, cols int
		fmt.Scan(&rows, &cols)
		matrix := InputMatrix(rows, cols)
		det, err := matrix.Determinant()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Determinant: %f\n", det)
		}
	case 6:
		fmt.Println("Enter dimensions for the matrix (rows and columns):")
		var rows, cols int
		fmt.Scan(&rows, &cols)
		matrix := InputMatrix(rows, cols)
		inverseMatrix, err := matrix.Inverse()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Inverse matrix:")
			fmt.Println(inverseMatrix)
		}
	case 7:
		fmt.Println("Enter dimensions for the matrix (rows and columns):")
		var rows, cols int
		fmt.Scan(&rows, &cols)
		matrix := InputMatrix(rows, cols)
		rowSorted := RowSortedMatrix{matrix}
		rowSorted.Sort()
		fmt.Println("Row-sorted matrix:")
		fmt.Println(rowSorted)
	case 8:
		fmt.Println("Enter dimensions for the matrix (rows and columns):")
		var rows, cols int
		fmt.Scan(&rows, &cols)
		matrix := InputMatrix(rows, cols)
		colSorted := ColSortedMatrix{matrix}
		colSorted.Sort()
		fmt.Println("Column-sorted matrix:")
		fmt.Println(colSorted)
	case 9:
		fmt.Println("Enter num of equations:")
		var rows int
		fmt.Scan(&rows)
		var cols = rows + 1
		matrix := InputMatrix(rows, cols)
		var result, err = matrix.SolveLinearEquations()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Solution for linear system:")
			fmt.Println(result)
		}
	default:
	}
}
