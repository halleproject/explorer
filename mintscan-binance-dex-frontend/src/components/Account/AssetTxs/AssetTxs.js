import React from "react";
import cn from "classnames/bind";
import styles from "./AssetTxs.scss";
import {_, empty} from "src/lib/scripts";
//  reduxy
import {useSelector} from "react-redux";
//  hooks
import {usePrevious} from "src/hooks";
//  components
import TxTable from "./TxTable";
import AssetsTable from "src/components/Account/AssetTxs/AssetsTable";

const cx = cn.bind(styles);

export default function({fetchAccountTxs = () => {}, balances = [], prices = [], txData = [], account = ""}) {
	const assets = useSelector(state => state.assets.assets);
	const [mappedAssets, setMappedAssets] = React.useState([]);
	const [selected, setSelected] = React.useState(true);
	const onClick = React.useCallback((e, bool) => {
		e.stopPropagation();
		e.preventDefault();
		setSelected(bool);
	}, []);

	React.useEffect(() => {
		if (!selected && empty(txData)) fetchAccountTxs();
	}, [fetchAccountTxs, selected, txData]);

	const prevBalances = usePrevious(balances);
	//  if balance has changed, just update that
	React.useEffect(() => {
		if (empty(prevBalances) || _.isEqual(prevBalances, balances) || empty(mappedAssets)) return;
		setMappedAssets(v => _.merge(v, balances));
		// console.log("merge new values");
	}, [balances, mappedAssets, prevBalances]);

	const txTable = React.useMemo(() => <TxTable txData={txData} account={account} />, [txData, account]);
	return React.useMemo(
		() => (
			<div className={cx("AssetTxs-wrapper")}>
				<div className={cx("Tabs")}>
					<div className={cx("Tab", !selected ? "selected" : undefined)} onClick={e => onClick(e, false)}>
						Transactions
					</div>
				</div>
				<div className={cx("Card")}>
					<div className={cx(!selected ? undefined : "unselected")}>{txTable}</div>
				</div>
			</div>
		),
		[onClick, selected, txTable]
	);
}
