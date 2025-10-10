package Engine

type Board[T any] struct {
	Board_num_rows int
	Board_num_cols int
	BoardData      [][]T

	TileGridTileSizeInPixels     int
	TileGridWidthInPixels        int
	TileGridHeightInPixels       int
	TileGridTotalXOffsetInPixels int
	TileGridTotalYOffsetInPixels int
}

func (board *Board[T]) InitBoard(boardNumCols, boardNumRows, boardGridTileSizeInPixels int, centreOnScreen bool) {

	board.Board_num_cols = boardNumCols
	board.Board_num_rows = boardNumRows

	board.BoardData = make([][]T, board.Board_num_rows)
	for i := range board.Board_num_rows {
		board.BoardData[i] = make([]T, board.Board_num_cols)
	}

	board.TileGridTileSizeInPixels = boardGridTileSizeInPixels
	board.TileGridWidthInPixels = board.Board_num_cols * board.TileGridTileSizeInPixels
	board.TileGridHeightInPixels = board.Board_num_rows * board.TileGridTileSizeInPixels
	if centreOnScreen {
		board.TileGridTotalXOffsetInPixels = (320 / 2) - (board.TileGridWidthInPixels / 2)
		board.TileGridTotalYOffsetInPixels = (240 / 2) - (board.TileGridHeightInPixels / 2)
	} else {
		board.TileGridTotalXOffsetInPixels = 0
		board.TileGridTotalYOffsetInPixels = 0
	}
}

func (board *Board[T]) GetPositionOfTileOnScreen(boardTile Vector2Int) Vector2Int {
	return Vector2Int{X: (boardTile.X * board.TileGridTileSizeInPixels) + board.TileGridTotalXOffsetInPixels, Y: (boardTile.Y * board.TileGridTileSizeInPixels) + board.TileGridTotalYOffsetInPixels}
}
