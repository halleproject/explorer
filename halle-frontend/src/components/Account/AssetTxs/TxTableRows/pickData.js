import React from "react";
import {_, formatNumber, reduceString, refineAddress} from "src/lib/scripts";
import {divide,fixed} from "src/lib/Big";
import getTxType from "src/constants/getTxType";
import {NavLink} from "react-router-dom";
import consts from "src/constants/consts";

//  components
import Skeleton from "react-skeleton-loader";
import Decimal from "src/components/common/Decimal";
//  assets
import greenArrowSVG from "src/assets/common/transferarrow_gr.svg";
import redArrowSVG from "src/assets/common/transferarrow_rd.svg";

import txTypes from "src/constants/txTypes";
const {COSMOS, WEB3, DEX, TOKENS, MISC} = txTypes;

export default function(data, cx, cell, account) {
	switch (cell) {
		case "txType": {
			if (!_.isNil(data?.messages[0].type)) return <span className={cx("type")}>{getTxType(data?.messages[0].type)}</span>;
			return <Skeleton />;
		}
		case "address": {
			var txtype = data?.messages[0].type
			if (txtype === WEB3.SEND || txtype == COSMOS.SEND) {
				if (data.fromAddr === "") {
					return <span>Multi send</span>;
				}
				const senderIsAcc = data.from_address === account;
				const baseAddr = refineAddress(senderIsAcc ? data.to_address : data.from_address);
				return (
					<NavLink className={cx("NavLink")} to={`/account/${baseAddr}`}>
						<img src={senderIsAcc ? redArrowSVG : greenArrowSVG} alt={"arrow"} />
						<span>{reduceString(baseAddr, 6, 6)}</span>
					</NavLink>
				);
			} else {
				//  No clickable since it's the same address
				return (
					<div>
						<span>{reduceString(data.from_address, 6, 6)}</span>
					</div>
				);
			}
		}
		case "Value": {
			return (
				<>
					{data.messages[0].type == "cosmos-sdk/MsgSend" ? (
						<div className={cx("number-display")}>
							<Decimal value={divide(data.messages[0].value.amount[0].amount, consts.NUM.BASE_MULT)} fontSizeBase={13} />
						</div>
					) : (
						<div className={cx("number-display")}>
							<Decimal value={divide(data.messages[0].value.value, consts.NUM.BASE_MULT)} fontSizeBase={13} />
						</div>
					)}
				</>
			);
		}
		case "Currency": {
			return <span className={cx({BNB: data.txAsset === "BNB"}, "text")}>{data.txAsset}</span>;
		}
		default:
			return "DEFAULT";
	}
}
