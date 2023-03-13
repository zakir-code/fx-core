// SPDX-License-Identifier: Apache-2.0

pragma solidity ^0.8.0;

interface ICrossChain {
    function crossChain(
        address _token,
        string memory _receipt,
        uint256 _amount,
        uint256 _fee,
        bytes32 _target,
        string memory _memo
    ) external payable returns (bool);

    function cancelSendToExternal(
        string memory _chain,
        uint256 _txid
    ) external returns (bool);
}

interface IERC20 {
    function allowance(
        address owner,
        address spender
    ) external view returns (uint256);

    function approve(address spender, uint256 amount) external returns (bool);

    function transferFrom(
        address from,
        address to,
        uint256 amount
    ) external returns (bool);
}

contract CrossChain is ICrossChain {
    address private constant _crossChainAddress =
        address(0x0000000000000000000000000000000000001004);

    function crossChain(
        address _token,
        string memory _receipt,
        uint256 _amount,
        uint256 _fee,
        bytes32 _target,
        string memory _memo
    ) external payable virtual override returns (bool) {
        if (_token != address(0)) {
            IERC20(_token).transferFrom(
                msg.sender,
                address(this),
                _amount + _fee
            );
            IERC20(_token).approve(_crossChainAddress, _amount + _fee);
        }
        return _crossChain(_token, _receipt, _amount, _fee, _target, _memo);
    }

    function _crossChain(
        address _token,
        string memory _receipt,
        uint256 _amount,
        uint256 _fee,
        bytes32 _target,
        string memory _memo
    ) internal returns (bool) {
        if (_token != address(0)) {
            uint256 allowance = IERC20(_token).allowance(
                address(this),
                _crossChainAddress
            );
            require(
                allowance == _amount + _fee,
                "allowance not equal amount + fee"
            );
        } else {
            require(
                msg.value == _amount + _fee,
                "msg.value not equal amount + fee"
            );
        }

        (bool result, bytes memory data) = _crossChainAddress.call{
            value: msg.value
        }(Encode.crossChain(_token, _receipt, _amount, _fee, _target, _memo));
        Decode.ok(result, data, "cross-chain failed");
        return Decode.crossChain(data);
    }

    function cancelSendToExternal(
        string memory _chain,
        uint256 _txid
    ) external virtual override returns (bool) {
        return _cancelSendToExternal(_chain, _txid);
    }

    function _cancelSendToExternal(
        string memory _chain,
        uint256 _txid
    ) internal returns (bool) {
        (bool result, bytes memory data) = _crossChainAddress.call(
            Encode.cancelSendToExternal(_chain, _txid)
        );
        Decode.ok(result, data, "cancel send to external failed");
        return Decode.cancelSendToExternal(data);
    }
}

library Encode {
    function crossChain(
        address _token,
        string memory _receipt,
        uint256 _amount,
        uint256 _fee,
        bytes32 _target,
        string memory _memo
    ) internal pure returns (bytes memory) {
        return
            abi.encodeWithSignature(
                "crossChain(address,string,uint256,uint256,bytes32,string)",
                _token,
                _receipt,
                _amount,
                _fee,
                _target,
                _memo
            );
    }

    function cancelSendToExternal(
        string memory _chain,
        uint256 _txid
    ) internal pure returns (bytes memory) {
        return
            abi.encodeWithSignature(
                "cancelSendToExternal(string,uin256)",
                _chain,
                _txid
            );
    }
}

library Decode {
    function crossChain(bytes memory data) internal pure returns (bool) {
        bool result = abi.decode(data, (bool));
        return result;
    }

    function cancelSendToExternal(
        bytes memory data
    ) internal pure returns (bool) {
        bool result = abi.decode(data, (bool));
        return result;
    }

    function ok(
        bool _result,
        bytes memory _data,
        string memory _msg
    ) internal pure {
        if (!_result) {
            string memory errMsg = abi.decode(_data, (string));
            if (bytes(_msg).length < 1) {
                revert(errMsg);
            }
            revert(string(abi.encodePacked(_msg, ": ", errMsg)));
        }
    }
}
