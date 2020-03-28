package kml

import (
	"image/color"
)

// A GxAltitudeModeEnum is an altitudeModeEnumType.
type GxAltitudeModeEnum string

// GxAltitudeModeEnums.
const (
	GxAltitudeModeClampToGround      GxAltitudeModeEnum = "clampToGround"
	GxAltitudeModeRelativeToGround   GxAltitudeModeEnum = "relativeToGround"
	GxAltitudeModeAbsolute           GxAltitudeModeEnum = "absolute"
	GxAltitudeModeClampToSeaFloor    GxAltitudeModeEnum = "clampToSeaFloor"
	GxAltitudeModeRelativeToSeaFloor GxAltitudeModeEnum = "relativeToSeaFloor"
)

// A GxFlyToModeEnum is a flyToModeEnumType.
type GxFlyToModeEnum string

// GxFlyToModeEnums.
const (
	GxFlyToModeBounce GxFlyToModeEnum = "bounce"
	GxFlyToModeSmooth GxFlyToModeEnum = "smooth"
)

// A GxPlayModeEnum is a playModeEnumType.
type GxPlayModeEnum string

// GxPlayModeEnums.
const (
	GxPlayModePause GxPlayModeEnum = "pause"
)

// GxAltitudeMode returns a new altitudeMode element.
func GxAltitudeMode(value GxAltitudeModeEnum) *SimpleElement {
	return newSEString("gx:altitudeMode", string(value))
}

// GxAltitudeOffset returns a new altitudeOffset element.
func GxAltitudeOffset(value float64) *SimpleElement {
	return newSEFloat("gx:altitudeOffset", value)
}

// GxBalloonVisibility returns a new balloonVisibility element.
func GxBalloonVisibility(value bool) *SimpleElement {
	return newSEBool("gx:balloonVisibility", value)
}

// GxDelayedStart returns a new delayedStart element.
func GxDelayedStart(value float64) *SimpleElement {
	return newSEFloat("gx:delayedStart", value)
}

// GxDrawOrder returns a new drawOrder element.
func GxDrawOrder(value int) *SimpleElement {
	return newSEInt("gx:drawOrder", value)
}

// GxDuration returns a new duration element.
func GxDuration(value float64) *SimpleElement {
	return newSEFloat("gx:duration", value)
}

// GxFlyToMode returns a new flyToMode element.
func GxFlyToMode(value GxFlyToModeEnum) *SimpleElement {
	return newSEString("gx:flyToMode", string(value))
}

// GxHorizFOV returns a new horizFov element.
func GxHorizFOV(value float64) *SimpleElement {
	return newSEFloat("gx:horizFov", value)
}

// GxInterpolate returns a new interpolate element.
func GxInterpolate(value bool) *SimpleElement {
	return newSEBool("gx:interpolate", value)
}

// GxLabelVisibility returns a new labelVisibility element.
func GxLabelVisibility(value bool) *SimpleElement {
	return newSEBool("gx:labelVisibility", value)
}

// GxOuterColor returns a new outerColor element.
func GxOuterColor(value color.Color) *SimpleElement {
	return newSEColor("gx:outerColor", value)
}

// GxOuterWidth returns a new outerWidth element.
func GxOuterWidth(value float64) *SimpleElement {
	return newSEFloat("gx:outerWidth", value)
}

// GxPhysicalWidth returns a new physicalWidth element.
func GxPhysicalWidth(value float64) *SimpleElement {
	return newSEFloat("gx:physicalWidth", value)
}

// GxPlayMode returns a new playMode element.
func GxPlayMode(value GxPlayModeEnum) *SimpleElement {
	return newSEString("gx:playMode", string(value))
}

// GxRank returns a new rank element.
func GxRank(value float64) *SimpleElement {
	return newSEFloat("gx:rank", value)
}

// GxValue returns a new value element.
func GxValue(value string) *SimpleElement {
	return newSEString("gx:value", value)
}

// GxX returns a new x element.
func GxX(value int) *SimpleElement {
	return newSEInt("gx:x", value)
}

// GxY returns a new y element.
func GxY(value int) *SimpleElement {
	return newSEInt("gx:y", value)
}

// GxW returns a new w element.
func GxW(value int) *SimpleElement {
	return newSEInt("gx:w", value)
}

// GxH returns a new h element.
func GxH(value int) *SimpleElement {
	return newSEInt("gx:h", value)
}

// GxAbstractTourPrimitive returns a new AbstractTourPrimitive element.
func GxAbstractTourPrimitive(children ...Element) *CompoundElement {
	return newCE("gx:AbstractTourPrimitive", children)
}

// GxAnimatedUpdate returns a new AnimatedUpdate element.
func GxAnimatedUpdate(children ...Element) *CompoundElement {
	return newCE("gx:AnimatedUpdate", children)
}

// GxFlyTo returns a new FlyTo element.
func GxFlyTo(children ...Element) *CompoundElement {
	return newCE("gx:FlyTo", children)
}

// GxPlaylist returns a new Playlist element.
func GxPlaylist(children ...Element) *CompoundElement {
	return newCE("gx:Playlist", children)
}

// GxSoundCue returns a new SoundCue element.
func GxSoundCue(children ...Element) *CompoundElement {
	return newCE("gx:SoundCue", children)
}

// GxTour returns a new Tour element.
func GxTour(children ...Element) *CompoundElement {
	return newCE("gx:Tour", children)
}

// GxTimeStamp returns a new TimeStamp element.
func GxTimeStamp(children ...Element) *CompoundElement {
	return newCE("gx:TimeStamp", children)
}

// GxTimeSpan returns a new TimeSpan element.
func GxTimeSpan(children ...Element) *CompoundElement {
	return newCE("gx:TimeSpan", children)
}

// GxTourControl returns a new TourControl element.
func GxTourControl(children ...Element) *CompoundElement {
	return newCE("gx:TourControl", children)
}

// GxWait returns a new Wait element.
func GxWait(children ...Element) *CompoundElement {
	return newCE("gx:Wait", children)
}

// GxLatLonQuad returns a new LatLonQuad element.
func GxLatLonQuad(children ...Element) *CompoundElement {
	return newCE("gx:LatLonQuad", children)
}

// GxTrack returns a new Track element.
func GxTrack(children ...Element) *CompoundElement {
	return newCE("gx:Track", children)
}

// GxMultiTrack returns a new MultiTrack element.
func GxMultiTrack(children ...Element) *CompoundElement {
	return newCE("gx:MultiTrack", children)
}

// GxSimpleArrayData returns a new SimpleArrayData element.
func GxSimpleArrayData(children ...Element) *CompoundElement {
	return newCE("gx:SimpleArrayData", children)
}

// GxViewerOptions returns a new ViewerOptions element.
func GxViewerOptions(children ...Element) *CompoundElement {
	return newCE("gx:ViewerOptions", children)
}
