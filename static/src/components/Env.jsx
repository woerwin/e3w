import React from 'react'
import  Cookies  from 'js-cookie'
import { Radio, Button } from 'antd'
import { GetEnvs } from "./request"

let envStyle = {
    position: 'absolute',top: '10px',left: '5px', fontSize: 18
}

const Env = React.createClass({
    getInitialState() {
        return {
            env: "",
            envList: [],
        }
    },

    _getEnvs(result) {
        console.log(result)
        let ec = localStorage.env
        if (result.indexOf(ec) === -1) {
            ec = result[0]
            localStorage.env = ec
        }
        this.setState({
            envList: result,
            env: ec,
        })
    },

    _get() {
        GetEnvs(this._getEnvs)
    },

    _setEnv(e) {
        this.setState({ env: e.target.value })
        localStorage.env = e.target.value
        // console.log(localStorage.env, e.target.value)
        location.replace(location.origin)
    },

    componentDidMount(){
        this._get()
    },

    render() {
        return (
            <div style={envStyle}>
                <Radio.Group value={this.state.env} onChange={this._setEnv} layo>
                {
                    this.state.envList.map(function (env) {
                        return <Radio.Button value={env} style={{display:"block"}}>{env}</Radio.Button>
                    })
                }
                </Radio.Group>
            </div>
        )
    }
})

module.exports = Env