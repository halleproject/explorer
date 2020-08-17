import React from "react";
import cn from "classnames/bind";
import styles from "./Address.scss";
import {_, empty, formatNumber} from "src/lib/scripts";
import copy from "copy-to-clipboard";
//  components
import {divide, fixed, multiply, sumArray} from "src/lib/Big";
import {useSelector} from "react-redux";
import Decimal from "src/components/common/Decimal";
//  assets
import qrSVG from "src/assets/account/qr_code.svg";
import QrModal from "../QrModal";
import DisplayLongString from "src/components/common/DisplayLongString";

const cx = cn.bind(styles);
export default function Address({account = {}, prices = [], assetData = {}}) {
	const bnbPrice = useSelector(state => state.blockchain.status?.price);
	const [showModal, setShowModal] = React.useState(false);
	const [copyClick, setCopyClick] = React.useState(false);
	// console.log(bnbPrice)
	const total = React.useMemo(() => {
		if (empty(prices)) return;
		const {balances} = account;
		const {assetData} = assetData.result.value.sequence;
		//  dollars
		const totalAssets = _.map(balances, (v, i) => multiply(getTotalSum(v), prices[i]));
		const freeAssets = _.map(balances, (v, i) => multiply(getFreeSum(v), prices[i]));
		// const freeAssets = _.map(balances, (v, i) => multiply(getFreeSum(v), prices[i]));

		return [sumArray(totalAssets), sumArray(freeAssets),assetData];
	}, [account, prices,assetData]);

	const renderQr = React.useMemo(
		() => (
			<div className={cx("qr-wrapper")} onClick={() => setShowModal(v => !v)} onMouseEnter={() => setShowModal(true)}
					 onMouseLeave={() => setShowModal(false)}>
				<img src={qrSVG} alt={"qr-code"}/>
				<QrModal address={account} show={showModal}/>
			</div>
		),
		[account, showModal],
	);
	const renderAddress = React.useMemo(() => {
		let timeout;
		const onClick = () => {
			if (copyClick) return;
			setCopyClick(true);
			copy(account);
			clearTimeout(timeout);
			timeout = setTimeout(() => {
				setCopyClick(false);
			}, 1000);
		};
		return (
			<ul className={cx("AddressDisplay-wrapper")}>
				<li className={cx("label")}>Address</li>
				<li className={cx("value")}>
					<p>
						{account ? (
							<>
								<DisplayLongString inputString={account} displayThresh={6}/>
								{/*<span className={cx("front")}>{account.substr(0, ADDRESS_CUT_DISPLAY_LENGTH)}</span>*/}
								{/*<span className={cx("remove")}>{account.substr(ADDRESS_CUT_DISPLAY_LENGTH, ADDRESS_LENGTH - ADDRESS_CUT_DISPLAY_LENGTH * 2)}</span>*/}
								{/*{account.substr(ADDRESS_LENGTH - ADDRESS_CUT_DISPLAY_LENGTH, ADDRESS_CUT_DISPLAY_LENGTH)}*/}
							</>
						) : (
							"-"
						)}
					</p>
					<img alt={"copy"} onClick={onClick}/>
					<div className={cx("copiedArea")}>
						<div className={cx("copied", {visible: copyClick})}>copied!</div>
					</div>
				</li>
			</ul>
		);
	}, [account, copyClick]);

	const renderStats = React.useMemo(
		() => (
			<div className={cx("statistics-wrapper")}>
				<ul className={cx("total-wrapper")}>
					<li className={cx("value")}>
						<span className={cx("front")}>Number of </span>
						<span className={cx("remove")}>transactions</span>
					</li>
					{/*<li className={cx("dollars")}>*/}
					{/*	<span className={cx("currency")}>$</span>*/}
					{/*	{!_.isNil(total?.[0]) && !_.isNil(bnbPrice) ? (*/}
					{/*		// <span>{formatNumber(fixed(total?.[0] ? total?.[0] : 0, 2))}</span>*/}
					{/*		<Decimal value={formatNumber(fixed(total?.[0] ? total?.[0] : 0, 2))} fontSizeBase={23} decimalReduce={5} bottom={2} />*/}
					{/*	) : (*/}
					{/*		<span>-</span>*/}
					{/*	)}*/}
					{/*</li>*/}
				</ul>
				<ul className={cx("compare-wrapper")}>
					<li className={cx("flexIt")}>
						{assetData?.result?.value?.sequence}
						{/*{prices.result}*/}
					</li>
					{/*<li className={cx("compareBNB")}>*/}
					{/*	{!_.isNil(total?.[0]) ? (*/}
					{/*		<>*/}
					{/*			/!*<span>{formatNumber(divide(total?.[0], bnbPrice, 2))}</span>*!/*/}
					{/*			<Decimal value={formatNumber(divide(total?.[0], bnbPrice, 2))} fontSizeBase={23} decimalReduce={5} bottom={2} />*/}
					{/*			<span className={cx("BNB")}>BNB</span>*/}
					{/*		</>*/}
					{/*	) : (*/}
					{/*		<span>-</span>*/}
					{/*	)}*/}
					{/*</li>*/}
				</ul>
			</div>
		),
		[bnbPrice, total,assetData],
	);

	return React.useMemo(
		() => (
			<div className={cx("Address-wrapper")}>
				<div className={cx("qr-address-wrapper")}>
					{renderQr}
					{renderAddress}
				</div>
				{renderStats}
			</div>
		),
		[renderAddress, renderQr, renderStats],
	);
}

const getTotalSum = asset => {
	return sumArray([asset.free, asset.locked, asset.frozen]);
};

const getFreeSum = asset => {
	return asset.free;
};
