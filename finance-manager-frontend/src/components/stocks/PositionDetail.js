import { useOutletContext } from "react-router-dom";
import { forwardRef, useEffect, useState } from "react";
import { LineChart } from '@mui/x-charts/LineChart'
import { ArrowDropUp, ArrowDropDown } from '@mui/icons-material'
import { format, parseISO } from "date-fns";

const PositionDetail = forwardRef((props, ref) => {

    const { numberFormatOptions } = useOutletContext();

    return (
        <div className="content">
            <hr />
            <h2>{props.position.ticker}</h2>
            <div className="d-flex justify-content-between">
                <div className="flex-col col-md-4 p-2">
                    <h3>Position Details:</h3>
                    <br/>
                    {
                        props.position ?
                            <>
                                <div className="d-flex justify-content-between">
                                    <div className="flex-col">
                                        <h4>Net Change:</h4>
                                    </div>
                                    <div className="flex-col">

                                        <h4 className={props.position.delta > 0 ? "text-success" : "text-failure"}>
                                            {props.position.delta > 0 ? <ArrowDropUp/> : <ArrowDropDown/>}
                                            ${Intl.NumberFormat("en-US", numberFormatOptions).format(props.position ? Math.abs(props.position.delta) : 0)}
                                        </h4>
                                    </div>
                                </div>
                                <div className="d-flex justify-content-between">
                                    <div className="flex-col">
                                        <h4>Net Change (%):</h4>
                                    </div>
                                    <div className="flex-col">
                                        <h4 className={props.position.delta > 0 ? "text-success" : "text-failure"}>
                                            {props.position.delta > 0 ? <ArrowDropUp/> : <ArrowDropDown/>}
                                            {Intl.NumberFormat("en-US", numberFormatOptions).format(props.position ? Math.abs(props.position.deltaPercentage) : 0)}%
                                        </h4>
                                    </div>
                                </div>
                                <div className="d-flex justify-content-between">
                                    <div className="flex-col">
                                        <h4>Open:</h4>
                                    </div>
                                    <div className="flex-col">
                                        <h4>${Intl.NumberFormat("en-US", numberFormatOptions).format(props.position ? props.position.open : 0)}</h4>
                                    </div>
                                </div>
                                <div className="d-flex justify-content-between">
                                    <div className="flex-col">
                                        <h4>Close:</h4>
                                    </div>
                                    <div className="flex-col">
                                        <h4>${Intl.NumberFormat("en-US", numberFormatOptions).format(props.position ? props.position.close : 0)}</h4>
                                    </div>
                                </div>
                                <div className="d-flex justify-content-between">
                                    <div className="flex-col">
                                        <h4>High:</h4>
                                    </div>
                                    <div className="flex-col">
                                        <h4>${Intl.NumberFormat("en-US", numberFormatOptions).format(props.position ? props.position.high : 0)}</h4>
                                    </div>
                                </div>
                                <div className="d-flex justify-content-between">
                                    <div className="flex-col">
                                        <h4>Low:</h4>
                                    </div>
                                    <div className="flex-col">
                                        <h4>${Intl.NumberFormat("en-US", numberFormatOptions).format(props.position ? props.position.low : 0)}</h4>
                                    </div>
                                </div>
                                {props.portfolioSummary &&
                                    <>
                                        <hr />
                                        <div className="d-flex justify-content-between">
                                            <div className="flex-col">
                                                <h4>Quantity:</h4>
                                            </div>
                                            <div className="flex-col">
                                                <h4>{props.portfolioSummary.quantity}</h4>
                                            </div>
                                        </div>
                                        <div className="d-flex justify-content-between">
                                            <div className="flex-col">
                                                <h4>Current Value:</h4>
                                            </div>
                                            <div className="flex-col">
                                                <h4>${Intl.NumberFormat("en-US", numberFormatOptions).format(props.portfolioSummary.value)}</h4>
                                            </div>
                                        </div>
                                    </>
                                }

                            </>

                            :
                            <h5>Failed to load Portfolio Data</h5>
                    }
                </div>
                <div className="flex-col col-md-8 p-2">
                    <h3>History:</h3>
                    {props.position.count > 0 ?
                        <LineChart
                            series={[
                                {
                                    data: props.position.values.map((v) => (v.open)),
                                    showMark: false,
                                    label: "open",
                                    color: "purple",
                                    id: "openData"
                                },
                                {
                                    data: props.position.values.map((v) => (v.close)),
                                    showMark: false,
                                    label: "close",
                                    color: props.position.values[0].close > props.position.values[props.position.count - 1].close ? "red" : "green",
                                    id: "closeData"
                                },
                                {
                                    data: props.position.values.map((v) => (v.high)),
                                    showMark: false,
                                    label: "high",
                                    color: "blue",
                                    id: "highData"
                                },
                                {
                                    data: props.position.values.map((v) => (v.low)),
                                    showMark: false,
                                    color: "orange",
                                    label: "low",
                                    id: "lowData"
                                },
                            ]}
                            xAxis={[{ scaleType: 'point', data: props.position.values.map((v) => format(parseISO(v.date), 'MMM do yyyy')) }]}
                            height={375}
                            slotProps={{
                                legend: { hidden: true }
                            }}
                            sx={{
                                '.MuiLineElement-series-openData': {
                                    stroke: 'none',
                                },
                                '.MuiLineElement-series-highData': {
                                    stroke: 'none',
                                },
                                '.MuiLineElement-series-lowData': {
                                    stroke: 'none',
                                }
                            }}
                        />
                        :
                        <h5>Data Not Available</h5>

                    }
                </div>
            </div>
        </div>

    );
});

export default PositionDetail;